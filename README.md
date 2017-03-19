## Metricsd

Reads OS metrics and sends data to AWS CloudWatch Metrics.
 
Metricsd gathers OS metrics for AWS CloudWatch. You can install it as a systemd process. 

Configuration
####  /etc/metricsd.conf 
```conf
# ------------------------------------------------------------
# env bool
#     Used to specify the environment: prod, dev, qa, staging, etc.
#     This gets used as a dimension that is sent to repeaters
#
# debug bool
#     true sets logging level to debug
#
# local bool
#     Used to ignore local ec2 meta-data, used for development only.
# ------------------------------------------------------------
env="dev"
#debug = true
#local = true

# ------------------------------------------------------------
# interval_seconds int
#     Defaults to 30 seconds
# ------------------------------------------------------------
#interval_seconds = 15

# ------------------------------------------------------------
# server_role
#     Used to specify the role of the AMI instance.
#     examples: dcos-master consul-master dcos-agent cassandra-node
# ------------------------------------------------------------
server_role = "dcos-master"

# ------------------------------------------------------------
# repeaters []string
#     aws, logger
#
# gatherers []string
#     disk cpu free nodetool
# ------------------------------------------------------------
repeaters = ["aws"]
gatherers = ["disk", "cpu", "free", "nodetool"]

# ------------------------------------------------------------
# aws_Region string
#     If not set, uses aws current region for this instance.
#     Used for testing only???
#
# ec2_instance_id string
#     If not set, uses aws instance id for this instance
#     Used for testing only???
#
# namespace string
#     Used to specify the top level namespace in cloudwatch.
# ------------------------------------------------------------
aws_region = "us-west-1"
ec2_instance_id = "my-fake-instanceid"
namespace="Cassandra Cluster"

# ------------------------------------------------------------
# disk_command string
#     default: /usr/bin/df
#     darwin example:  /bin/df
#
# disk_file_systems []string
#     What FileSystems to include. defaults to /dev/*
#
# disk_fields []string
#     what fields to output
#     fields: total        - number of 1K bytes on the disk
#             used         - number of 1K bytes used on the disk
#             available    - number of 1K bytes available on the disk
#             usedpct      - percentage of bytes used on the disk (calculated)
#             availablepct - percentage of bytes available on the disk (calculated)
#             capacitypct  - percentage of bytes available on the disk (reported)
#             mount        - where FileSystem is mounted
#     default: ["availablepct"]
# ------------------------------------------------------------
#disk_command = "/usr/mybin/df"
#disk_file_systems = ["/dev/*", "udev"]
#disk_fields = ["total", "used", "available", "usedpct", "availablepct", "mount"]

# ------------------------------------------------------------
# cpu_proc_stat string
#     Used to specify /proc/stat or absolute file
#     default:        /proc/stat
#     darwin example: /home/batman/gospace/src/github.com/cloudurable/metricsd/test/test-data/proc/stat
#
# cpu_report_zeros bool
#     default: false
# ------------------------------------------------------------
#cpu_proc_stat = "/proc/stat"
#cpu_report_zeros = true

# ------------------------------------------------------------
# free_command string
#     default:        /usr/bin/free
#     darwin example: /usr/local/bin/free
# ------------------------------------------------------------
#free_command = "/usr/mybin/free"

# ------------------------------------------------------------
# nodetool_command string
#   default:        /usr/bin/nodetool
#   darwin example: /usr/local/bin/nodetool
#
# nodetool_functions []string    : required when nodetool will run
#    specify nodetool_functions as array
#      functions: cfstats compactionstats gcstats netstats tpstats getlogginglevels proxyhistograms listsnapshots, statuses
# ------------------------------------------------------------
#nodetool_command = "/usr/mybin/nodetool"
#nodetool_functions = ["tpstats", "gcstats"]
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
