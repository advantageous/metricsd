## Metricsd

Reads OS metrics and sends data to AWS CloudWatch Metrics.
 
Metricsd gathers OS metrics for AWS CloudWatch. You can install it as a systemd process. 

Configuration
####  /etc/metricsd.conf 
```conf
# ------------------------------------------------------------
# AWS Region         string        `hcl:"aws_region"`
# If not set, uses aws current region for this instance.
# Used for testing only.
# ------------------------------------------------------------
# aws_region = "us-west-1"

# ------------------------------------------------------------
# EC2InstanceId     string        `hcl:"ec2_instance_id"`
# If not set, uses aws instance id for this instance
# Used for testing only.
# ------------------------------------------------------------
# ec2_instance_id = "i-my-fake-instanceid"

# ------------------------------------------------------------
# Debug             bool          `hcl:"debug"`
# Used for testing and debugging
# ------------------------------------------------------------
debug = false

# ------------------------------------------------------------
# Local             bool          `hcl:"local"`
# Used to ingore local ec2 meta-data, used for development only.
# ------------------------------------------------------------
# local = true

# ------------------------------------------------------------
# TimePeriodSeconds time.Duration `hcl:"interval_seconds"`
# Defaults to 30 seconds, how often metrics are collected.
# ------------------------------------------------------------
interval_seconds = 10

# ------------------------------------------------------------
# Used to specify the environment: prod, dev, qa, staging, etc.
# This gets used as a dimension that is sent to cloudwatch. 
# ------------------------------------------------------------
env="dev"

# ------------------------------------------------------------
# Used to specify the top level namespace in cloudwatch.
# ------------------------------------------------------------
namespace="Cassandra Cluster"

# ------------------------------------------------------------
# Used to specify the role of the AMI instance.
# Gets used as a dimension.
# e.g., dcos-master, consul-master, dcos-agent, cassandra-node, etc.
# ------------------------------------------------------------
server_role="dcos-master"

# ------------------------------------------------------------
# used to specify Disk gatherer properties and if it runs
#
# default disk_command: /usr/bin/df
# darwin disk_command:  /bin/df
#
# default disk_includes:    "/dev/*"               includes all FileSystems starting with /dev/
# example specific include: "/dev/sda5"            includes only /dev/sda5
# example multiple include: "/dev/sda5 /other/*"   includes /dev/sda5 and all starting with /other/* evaluated in order in list
# ------------------------------------------------------------
disk_gather = true
disk_command = "df"
#disk_includes = "/dev/*"
#disk_includes =
# ------------------------------------------------------------
# used to specify Cpu gatherer properties and if it runs
#
# default cpu_proc_stat: /proc/stat
# darwin cpu_proc_stat:  /home/rickhigh/gospace/src/github.com/cloudurable/metricsd/metric/test-data/proc/stat
# ------------------------------------------------------------
cpu_gather = true
#cpu_proc_stat = "/proc/stat"

# ------------------------------------------------------------
# used to specify Free Ram gatherer properties and if it runs
#
# default free_command: /usr/bin/free
# darwin free_command:  /usr/local/bin/free
# ------------------------------------------------------------
free_gather = true
#free_command = "/usr/bin/free"

# ------------------------------------------------------------
# used to specify Nodetool gatherer properties and if it runs
#
# default nodetool_command: /usr/bin/nodetool
# darwin nodetool_command:  /usr/local/bin/nodetool
#
# specify nodetool_functions with space delimted string
# all nodetool_functions: "cfstats compactionstats gcstats netstats tpstats getlogginglevels"
# ------------------------------------------------------------
nodetool_gather = true
#nodetool_command = "/usr/bin/nodetool"
nodetool_functions = "gcstats"
```


## Installing as a service

If you are using systemd you should install this as a service. 

#### /etc/systemd/system/metricsd.service
```conf

[Unit]
Description=metricsd
Wants=basic.target
After=basic.target network.target

[Service]
User=centos
Group=centos
ExecStart=/usr/bin/metricsd
KillMode=process
Restart=on-failure
RestartSec=42s


[Install]
WantedBy=multi-user.target


```
Copy the binary to `/usr/bin/metricsd`.
Copy the config to `/etc/metricsd.conf`.
You can specify a different conf location by using `/usr/bin/metricsd -conf /foo/bar/myconf.conf`.

#### Installing 
```sh

$ sudo cp metricsd_linux /usr/bin/metricsd 
$ sudo systemctl stop  metricsd.service
$ sudo systemctl enable  metricsd.service
$ sudo systemctl start  metricsd.service
$ sudo systemctl status  metricsd.service
● metricsd.service - metricsd
   Loaded: loaded (/etc/systemd/system/metricsd.service; enabled; vendor preset: disabled)
   Active: active (running) since Wed 2016-12-21 20:19:59 UTC; 8s ago
 Main PID: 718 (metricsd)
   CGroup: /system.slice/metricsd.service
           └─718 /usr/bin/metricsd

Dec 21 20:19:59 ip-172-31-29-173 systemd[1]: Started metricsd.
Dec 21 20:19:59 ip-172-31-29-173 systemd[1]: Starting metricsd...
Dec 21 20:19:59 ip-172-31-29-173 metricsd[718]: INFO     : [main] - 2016/12/21 20:19:59 config.go:29: Loading config /et....conf
Dec 21 20:19:59 ip-172-31-29-173 metricsd[718]: INFO     : [main] - 2016/12/21 20:19:59 config.go:45: Loading log...
```

There are full example packer install scripts under bin/packer/packer_ec2.json.
The best doc is a working example. 

## Metrics

### CPU metrics
* `softIrqCnt` - count of soft interrupts for the last period
* `intrCnt` - count of interrupts for the last period
* `ctxtCnt` - count of context switches for the last period
* `processesStrtCnt` - count of processes started for the last period
* `GuestJif` - jiffies spent in guest mode for last time period
* `UsrJif` - jiffies spent in usr mode for last time period
* `IdleJif` - jiffies spent in usr mode for last time period
* `IowaitJif` - jiffies spent handling IO for last time period
* `IrqJif` - jiffies spent handling interrupts for last time period
* `GuestniceJif` - guest nice mode
* `StealJif` - time stolen by noisy neighbors for last time period
* `SysJif` - jiffies spent doing OS stuff like system calls in last time period
* `SoftIrqJif` - jiffies spent handling soft IRQs in the last time period
* `procsRunning` - count of processes currently running
* `procsBlocked` - count of processes currently blocked (could be for IO or just waiting to get CPU time)


### Disk metrics
* `dUVol<VOLUME_NAME>AvailPer` - percentage of disk space left (per volume)

### Mem metrics
* `mFreeLvl` - free memory in kilobytes
* `mUsedLvl` - used memory in kilobytes
* `mSharedLvl` - shared memory in kilobytes
* `mBufLvl` - memory used by IO buffers in kilobytes
* `mAvailableLvl` - memory available in kilobytes
* `mFreePer` - percentage of memory free
* `mUsedPer` - percentage of memory used

If swapping is enabled (which is unlikely), then you will get the above with `mSwpX` instead of `mX`.


### Nodetool gcstats
* `gcInterval` - Interval (ms)
* `gcMaxElapsed` - Max GC Elapsed (ms)
* `gcTotalElapsed` - Total GC Elapsed (ms)
* `gcStdevElapsed` - Stdev GC Elapsed (ms)
* `gcReclaimed` - GC Reclaimed (MB)
* `gcCollections` - Collections
* `gcDirectMemoryBytes` - Direct Memory Bytes

### Nodetool netstats
* `nsMode` - mode of the node
 * STARTING = 0
 * NORMAL = 1
 * JOINING = 2
 * LEAVING = 3
 * DECOMMISSIONED = 4
 * MOVING = 5
 * DRAINING = 6
 * DRAINED = 7
 * OTHER = 99

* `nsRrAttempted` - Read Repair Attempted
* `nsRrBlocking` - Read Repair Mismatch (Blocking)
* `nsRrBackground` - Read Repair Mismatch (Background)

All message counts can be >= 0, or n/a (-125) or error (-127)
* `nsPoolLargeMsgsActive` - Large messages Active
* `nsPoolLargeMsgsPending` - Large messages Pending
* `nsPoolLargeMsgsCompleted` - Large messages Completed
* `nsPoolLargeMsgsDropped` - Large messages Dropped
* `nsPoolSmallMsgsActive` - Small messages Active
* `nsPoolSmallMsgsPending` - Small messages Pending
* `nsPoolSmallMsgsCompleted` - Small messages Completed
* `nsPoolSmallMsgsDropped` - Small messages Dropped
* `nsPoolGossipMsgsActive` - Gossip messages Active
* `nsPoolGossipMsgsPending` - Gossip messages Pending
* `nsPoolGossipMsgsCompleted` - Gossip messages Completed
* `nsPoolGossipMsgsDropped` - Gossip messages Dropped
