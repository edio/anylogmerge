# anylogmerge

anylogmerge is a tool that allows to merge interleaved logs.

It allows you to define regular expression that will determine how log entries are compared. For example consider log file of the following structure

    DEBUG 2013/06/06 23:16:47:598 [UserService] Attempt to login. User=JohnDoe
    DEBUG 2013/06/06 23:16:47:604 [PreparedStatement] ==>  Executing: select * from USERS where login = ?
    DEBUG 2013/06/06 23:16:47:605 [PreparedStatement] ==> Parameters: JohnDoe(String) 
    ERROR 2013/06/06 23:18:54:515 [UserService] Error while processing login request. No such user found
    com.example.myapp.UserNotFoundException 
        at com.example.myapp.service.UserServiced.login(UserService.java:143) ~[myapp.war:na]
        at sun.reflect.GeneratedMethodAccessor912.invoke(Unknown Source) ~[na:na]
        at sun.reflect.DelegatingMethodAccessorImpl.invoke(DelegatingMethodAccessorImpl.java:25) ~[na:1.6.0_26]
        at java.lang.reflect.Method.invoke(Method.java:597) ~[na:1.6.0_26]
        at org.springframework.aop.support.AopUtils.invokeJoinpointUsingReflection(AopUtils.java:318) [app.war:3.1.0.RELEASE]
        at org.springframework.aop.framework.JdkDynamicAopProxy.invoke(JdkDynamicAopProxy.java:196) [app.war:3.1.0.RELEASE]
        at $Proxy305.login(Unknown Source) ~[na:na]
    INFO 2013/06/06 23:26:59:746 [CartService] Cart[id=1234] checked out. User=NeoTheOne

To merge files of such structure properly one have to extract timestamps and compare log entries by timestamp while merging. anylogmerge allows this. For instance to merge such logs

    anylogmerge -s "^\w+\ (\S+)\ (\S+).*$" -o merged.log node1/app.log node2/app.log node3/app.log

This will make anylogmerge to compare log entries by substrings that match capturing groups in provided regular expression. In this particular case log entries will be compared by values like

    2013/06/0623:16:47.598

## Operation modes

anylogmerge supports 3 comparing modes for log entries while merging

* compare log entries with each other as is
* compare by key, which is defined by column specification (explained below)
* compare by key, defined by regex

### Simple mode

The first one, the simplest, is enabled by default. It is suitable for logs where each line starts with a sortable field 

    2012/12/31 23:59:59 Ding!
    2013/01/01 00:00:00 Happy New Year!

Usage:

    anylogmerge log1.log log2.log .. logN.log

### Columns mode

The second one allows to use substring of a log entry as a comparing key. Substring is defined as comma-separated pairs of indexes

    :3,5:10,15:

Expression above tells anylogmerge to extract substrings at 0 to 3, 5 to 10 and 15 to len(string) positions and concatenate them into a key.
This one is usefull, when log entries do not start with a sortable field but have distinct columns

    DEBUG 2012/12/31 23:59:59 Ding!
    INFO  2013/01/01 00:00:00 Happy New Year!

Usage

    anylogmerge -c "6:16,17:25" log1.log log2.log .. logN.log

If launched in that way against examplary log, key would be *2012/12/3123:59:59*.

Though this mode is enough for many cases, it may produce wrong output for multiline logs

### Regex mode

The last and most powerfull mode allows to define comparing key with regex.

    DEBUG 2012/12/31 23:59:59 Ding!
    INFO 2013/01/01 00:00:00 Happy New Year!
    Best wishes and goodluck
    DEBUG 2013/01/01 00:00:01 The first second of 2013 year

In the log above, date is not placed at fixed column but drifts by 1 place to the left for INFO log entries.
Moreover, some log entries can take several lines and we want such entries appear in resulting log unsplitted.

So we have to define regular expression to extract date from the log

    ^[A-Z]+\ ([0-9\/]+)\ ([0-9\:]+).*$

This regular expression has 2 capturing groups that define comparing key. Thus generic log entry with the expression above will be compared by key *2012/12/3123:59:59*.
What is also great, is that unmatched log entries will have empty string key. Empty string is lexicographically less than the other strings, so unmatched log entries are passed to output in the first place.
This is how anylogmerge handles multiline logs.

Usage

    anylogmerge -s "^\[A-Z]+\ ([0-9\/]+)\ ([0-9\:]+).*$" log1.log log2.log .. logN.log

# Performance

Performance of the anylogmerge naturally depends on the mode of operation.

For example while merging 2 ~50MB files in different modes, the following results achieved on Core2Duo E6550 with an old 7200 RPM HDD

    $ time anylogmerge -v -f -o ~/tmp/logs/merged.log ~/tmp/logs/app-all.*.log                         
    Merging files [/home/edio/tmp/logs/app-all.node-155.log /home/edio/tmp/logs/app-all.node-156.log]
    anylogmerge -v -f -o ~/tmp/logs/merged.log   0.92s user 1.49s system 83% cpu 3.010 total

    $ time anylogmerge -c "6:30" -v -f -o ~/tmp/logs/merged.log ~/tmp/logs/app-all.*.log                         
    Merging files [/home/edio/tmp/logs/app-all.node-155.log /home/edio/tmp/logs/app-all.node-156.log]
    anylogmerge -c "6:30" -v -f -o ~/tmp/logs/merged.log   1.02s user 1.50s system 84% cpu 2.972 total

    $ time anylogmerge -s "^\w+\ (\S+)\ (\S+).*$" -v -f -o ~/tmp/logs/merged.log ~/tmp/logs/app-all.*.log                         
    Merging files [/home/edio/tmp/logs/app-all.node-155.log /home/edio/tmp/logs/app-all.node-156.log]
    Matched key example '^\w+\ (\S+)\ (\S+).*$' : '2013/06/0623:09:46:192'
    anylogmerge -s "^\w+\ (\S+)\ (\S+).*$" -v -f -o ~/tmp/logs/merged.log   18.15s user 1.15s system 100% cpu 19.298 total

    $ ls -l ~/tmp/logs
    total 207864
    -rw-r--r-- 1 edio users 106419456 Jun 16 19:29 merged.log
    -rw-r--r-- 1 edio users  59278157 Jun  9 03:00 app-all.node-155.log
    -rw-r--r-- 1 edio users  47141299 Jun  9 03:00 app-all.node-156.log

So for simple and columns mode, performance is bounded by IO. Merging to /dev/null takes less than a second in these modes.
Regex mode is much slower, and merging the same files takes about 20 seconds to complete.

# Building

anylogmerge is written in Go, so you need one to build the tool.
I'll provide some building script to simplify building for those, who are not familiar with Go.

For now, only brief building instructions.

1. git clone project to local directory
2. cd to directory where you have cloned the project
3. build

    $ export GOPATH=`pwd`
    $ go install logmerge
    $ go build src/anylogmerge.go

# Why Go

Because I'm trying to learn Go.

# License

GPLv3

