---
title: 使用inotfy 扩展 让swoole热加载
date: 2019-3-8
categories: 
- php
- swoole
tags:
- inotfy
- swoole
---

## inotify 
### inotify  介绍
>* inotify是**Linux内核提供的一组系统调用**，它可以监控文件系统操作
，比如文件或者目录的创建、读取、写入、权限修改和删除等。
 inotify使用也很简单，使用**inotify_init**创建一个句柄，然后通过**inotify_add_watch/inotify_rm_watch**增加/删除对文件和目录的监听。
  PHP中提供了inotify扩展，支持了inotify系统调用。inotify本身也是一个文件描述符，可以加入到事件循环中，配合使用swoole扩展
  ，就可以**异步非阻塞地实时监听文件/目录变化**
 
### php 安装 inotify 
* 1、可以使用 pecl install inotify
* 2、（https://pecl.php.net/get/inotify-2.0.0.tgz ）编译安装步骤跟之前一样，自行安装
 
### 实战


#### Watch.php
```php
<?php

namespace Swoole\ToolKit;

class NotFound extends \Exception
{
}

class Watch
{
    /**
     * @var resource
     */
    protected $inotify;
    protected $pid;
    protected $reloadFileTypes = array('.php' => true);
    protected $watchFiles = array();
    protected $afterNSeconds = 10;
    /**
     * 正在reload
     */
    protected $reloading = false;
    protected $events;
    /**
     * 根目录
     * @var array
     */
    protected $rootDirs = array();

    function putLog($log)
    {
        $_log = "[" . date('Y-m-d H:i:s') . "]\t" . $log . "\n";
        echo $_log;
    }

    /**
     * @param $serverPid
     * @throws NotFound
     */
    function __construct($serverPid)
    {
        $this->pid = $serverPid;
        if (posix_kill($serverPid, 0) === false) {
            throw new NotFound("Process#$serverPid not found.");
        }
        $this->inotify = inotify_init();
        $this->events = IN_MODIFY | IN_DELETE | IN_CREATE | IN_MOVE;
        swoole_event_add($this->inotify, function ($ifd) {
            $events = inotify_read($this->inotify);
            if (!$events) {
                return;
            }
            //var_dump($events);
            foreach ($events as $ev) {
                if ($ev['mask'] == IN_IGNORED) {
                    continue;
                } else if ($ev['mask'] == IN_CREATE or $ev['mask'] == IN_DELETE or $ev['mask'] == IN_MODIFY or $ev['mask'] == IN_MOVED_TO or $ev['mask'] == IN_MOVED_FROM) {
                    $fileType = strrchr($ev['name'], '.');
                    //非重启类型
                    if (!isset($this->reloadFileTypes[$fileType])) {
                        continue;
                    }
                }
                //正在reload，不再接受任何事件，冻结10秒
                if (!$this->reloading) {
                    $this->putLog("after 10 seconds reload the server");
                    //有事件发生了，进行重启
                    swoole_timer_after($this->afterNSeconds * 1000, array($this, 'reload'));
                    $this->reloading = true;
                }
            }
        });
    }

    function reload()
    {
        $this->putLog("reloading");
        //向主进程发送信号
        posix_kill($this->pid, SIGUSR1);
        //清理所有监听
        $this->clearWatch();
        //重新监听
        foreach ($this->rootDirs as $root) {
            $this->watch($root);
        }
        //继续进行reload
        $this->reloading = false;
    }

    /**
     * 添加文件类型
     * @param $type
     */
    function addFileType($type)
    {
        $type = trim($type, '.');
        $this->reloadFileTypes['.' . $type] = true;
    }

    /**
     * 添加事件
     * @param $inotifyEvent
     */
    function addEvent($inotifyEvent)
    {
        $this->events |= $inotifyEvent;
    }

    /**
     * 清理所有inotify监听
     */
    function clearWatch()
    {
        foreach ($this->watchFiles as $wd) {
            inotify_rm_watch($this->inotify, $wd);
        }
        $this->watchFiles = array();
    }

    /**
     * @param $dir
     * @param bool $root
     * @return bool
     * @throws NotFound
     */
    function watch($dir, $root = true)
    {
        //目录不存在
        if (!is_dir($dir)) {
            throw new NotFound("[$dir] is not a directory.");
        }
        //避免重复监听
        if (isset($this->watchFiles[$dir])) {
            return false;
        }
        //根目录
        if ($root) {
            $this->rootDirs[] = $dir;
        }
        $wd = inotify_add_watch($this->inotify, $dir, $this->events);
        $this->watchFiles[$dir] = $wd;
        $files = scandir($dir);
        foreach ($files as $f) {
            if ($f == '.' or $f == '..') {
                continue;
            }
            $path = $dir . '/' . $f;
            //递归目录
            if (is_dir($path)) {
                $this->watch($path, false);
            }
            //检测文件类型
            $fileType = strrchr($f, '.');
            if (isset($this->reloadFileTypes[$fileType])) {
                $wd = inotify_add_watch($this->inotify, $path, $this->events);
                $this->watchFiles[$path] = $wd;
            }
        }
        return true;
    }

    function run()
    {
        swoole_event_wait();
    }
}

```
