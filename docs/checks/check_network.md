# check_network

Checks network interface status.

- [Examples](#examples)
- [Argument Defaults](#argument-defaults)
- [Metrics](#metrics)

## Examples

### **Default check**

    check_network
    OK: lo >11063848400 <11063848400 bps, eth0 >31774542210 <121709796043 bps

### Example using **NRPE** and **Naemon**

Naemon Config

    define command{
        command_name    check_nrpe
        command_line    $USER1$/check_nrpe -2 -H $HOSTADDRESS$ -n -c $ARG1$ -a $ARG2$
    }

    define service {
            host_name               testhost
            service_description     check_network_testhost
            check_command           check_nrpe!check_network!'crit=none'
    }

Return

    OK: lo >11063848400 <11063848400 bps, eth0 >31774542210 <121709796043 bps

## Argument Defaults

| Argument | Default Value |
| --- | --- |
top-syntax | %(status): %(list) |
ok-syntax | %(status): %(list) |
detail-syntax | %(name) >%(sent) <%(received) bps |

## Metrics

#### **Check specific metrics**

| Metric | Description |
| --- | --- |
| MAC | The MAC address |
| enabled | True if the network interface is enabled |
| name | Name of the interface |
| net_connection_id | Network connection id |
| received | Bytes received per second |
| sent | Bytes sent per second |
| speed | Network interface speed? |