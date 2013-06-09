anylogmerge
===========

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
        at org.springframework.aop.support.AopUtils.invokeJoinpointUsingReflection(AopUtils.java:318) [tcs.war:3.1.0.RELEASE]
        at org.springframework.aop.framework.JdkDynamicAopProxy.invoke(JdkDynamicAopProxy.java:196) [tcs.war:3.1.0.RELEASE]
        at $Proxy305.login(Unknown Source) ~[na:na]
    INFO 2013/06/06 23:26:59:746 [CartService] Cart[id=1234] checked out. User=NeoTheOne

To merge files of such structure properly one have to extract timestamps and compare log entries by timestamp while merging. anylogmerge allows this. For instance to merge such logs

    anylogmerge -s "^\w+\ (\S+)\ (\S+).*$" -o merged.log node1/app.log node2/app.log node3/app.log

This will make anylogmerge to compare log entries by substrings that match capturing groups in provided regular expression. In this particular case log entries will be compared by values like

    2013/06/0623:16:47.598

Performance
-----------

I'm just learning Go and tool doesn't aim to perform like insane. I'm sure there is plenty of room for improvement. 
On Core2Duo E6550 with an old 7200 RPM HDD it takes less than 20 seconds to merge 2 ~50MB log files.

    $ ./anylogmerge -s "^\w+\ (\S+)\ (\S+).*$" -v -f -o ~/tmp/logs/merged.log ~/tmp/logs/*.log                         
    Merging files [/home/edio/tmp/logs/tcs-all.155.log /home/edio/tmp/logs/tcs-all.156.log]
    Matched key example '^\w+\ (\S+)\ (\S+).*$' : '2013/06/0623:09:46:192'
    ./anylogmerge -s "^\w+\ (\S+)\ (\S+).*$" -v -f -o ~/tmp/logs/merged.log   18.15s user 1.15s system 100% cpu 19.298 total
    $ ls -l ~/tmp/logs
    total 207864
    -rw-r--r-- 1 edio users 106419456 Jun  9 03:02 merged.log
    -rw-r--r-- 1 edio users  59278157 Jun  9 03:00 tcs-all.155.log
    -rw-r--r-- 1 edio users  47141299 Jun  9 03:00 tcs-all.156.log


Building
--------

anylogmerge is written in Go, so you need one to build the tool.
I'll provide some building script to simplify building for those, who are not familiar with Go.



Why Go
------

Because I'm trying to learn Go.

