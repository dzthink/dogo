// Copyright 2017 dzthink
// license that can be found in the LICENSE file.

//provide a cross platform for some syscall api

// +build !windows

package dogo
import (
    "os"
    "syscall"
    "os/signal"
    "path"
    "io/ioutil"
    "strconv"
    "os/user"
)

func logProcessInfo(config *Config) {
    //记录进程id
    var(
        err error
        pid string
        uStr string
        gStr string
    )
    if pid, err = config.String("process.pid"); err == nil {
        pid, err = os.Getwd()
        if err == nil{
            pid = "/tmp/dogo.pid"
        } else {
            pid = pid + string(os.PathSeparator) + "dogo.pid"
        }
    }
    os.MkdirAll(path.Dir(pid), 0755)
    os.Create(pid)

    ioutil.WriteFile(pid, []byte(strconv.Itoa(os.Getpid())), 0755)

    //设置用户和组信息

    if uStr, err = config.String("process.user"); err == nil {
       u, e := user.Lookup(uStr)
       if e == nil {
           uid, _ := strconv.Atoi(u.Uid)
           syscall.Setuid(uid)
       }
    }
    if gStr, err = config.String("process.gid"); err == nil {
       g, e := user.LookupGroup(gStr)
       if e == nil {
           gid, _ := strconv.Atoi(g.Gid)
           syscall.Setgid(gid)
       }
    }
}

func processSignal() {
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL,syscall.SIGUSR1, syscall.SIGUSR2, os.Interrupt)
    for{
        msg := <-sigs
        switch msg {
        default:

        case syscall.SIGUSR1:
            //reload

        case syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM:
            //logger.Info("application stoping, signal[%v]", msg)
            //b.App.Stop()
            return
        }
    }
}




