
Sidecar nsm-monitor
============================

Specification
-------------
nsm-monitor is simple sidecar for monitoring connection events into client's POD. Also, it providing API for creating custom nsm-monitor (see below).  

Implementation details
---------------------------------

#### How to start to monitor into client's POD
For monitoring client's connection events you need to add into POD additional container. For example:
```
...
spec:
    spec:
      hostPID: true
      containers:
        #original container
        - name: vppagent-nsc
          image: networkservicemesh/vpp-test-common:latest
          imagePullPolicy: IfNotPresent
        #injected monitor
        - name: nsm-monitor
          image: networkservicemesh/nsm-monitor:lateest
          imagePullPolicy: IfNotPresent
...

```

#### How to create custom nsm-monitor
For creating custom nsm-monitor you need to implement next interface: 
```
type NSMMonitorHelper interface {
    //Runs on nsm-monitor connected
    Connected(map[string]*connection.Connection)
    //Runs on healing started
    Healing(conn *connection.Connection)
    //Gets custom network service configuration
    GetConfiguration() *common.NSConfiguration
    /Runs on restore failed, the error pass as the second parameter
    ProcessHealing(newConn *connection.Connection, e error)
    //Runs on invoked NSMMonitorApp.Stop()
    Stopped()
    //Returns is Jaeger needed
    IsEnableJaeger() bool
}
```
After that, you could use next code for your sidecar:
```
func main() {
    c := tools.NewOSSignalChannel()
    app := nsm_sidecars.NewNSMMonitorApp()
    app.SetHelper(newMyHepler())
    go app.Run(version)
    <-c
}
```

Example usage
------------------------
For an example of usage you could take a look at tests:

* TestNSMMonitorInit
* TestDeployNSMMonitor

References
----------

* [Spec: ResiliencyV2](https://github.com/networkservicemesh/networkservicemesh/issues/1331) scenario №9.
* [Soec: Spec: DNS Integration for NSM](https://github.com/networkservicemesh/networkservicemesh/issues/1224) scenario №1