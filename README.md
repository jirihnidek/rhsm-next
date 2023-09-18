Red Hat Subscription Management Service 2.0
===========================================

The Red Hat Subscription Management Service provides D-Bus API. It is implemented
in Go.

## Build instruction

```
go build
```

## Testing rhsm-next

Start `rhsm-next` using in Terminal #1:

```
sudo -u rhsm2 ./rhsm-next
```

Then you can call methods for each interface in Terminal #2

### com.redhat.RHSM2.Consumer

```
busctl --system call com.redhat.RHSM2 /com/redhat/RHSM2/Consumer \
   com.redhat.RHSM2.Consumer GetOrg s "en_US.utf-8"
```

```
busctl --system call com.redhat.RHSM2 /com/redhat/RHSM2/Consumer \
   com.redhat.RHSM2.Consumer GetUuid s "en_US.utf-8"
```

### com.redhat.RHSM2.Register

```
busctl --system call com.redhat.RHSM2 /com/redhat/RHSM2/Register \
    com.redhat.RHSM2.Register GetOrgs sss "admin" "admin" ""
```

```
busctl --system call com.redhat.RHSM2 /com/redhat/RHSM2/Register \
    com.redhat.RHSM2.Register RegisterWithUsername sssa{ss}s "donaldduck" "admin" "admin" 0 ""
```

```
busctl --system call com.redhat.RHSM2 /com/redhat/RHSM2/Register \
    com.redhat.RHSM2.Register RegisterWithActivationKey sasa{ss}s "donaldduck" 1 "awesome_os_pool" 0 ""
```