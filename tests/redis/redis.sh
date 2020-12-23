#!/bin/bash

cluster () {
    ipaddr=$(hostname -i)
    # 邪魔なファイルを削除。
    rm -f \
        /data/conf/r6379i.log \
        /data/conf/r6380i.log \
        /data/conf/r6381i.log \
        /data/conf/r6382i.log \
        /data/conf/r6383i.log \
        /data/conf/r6384i.log \
        /data/conf/nodes.6379.conf \
        /data/conf/nodes.6380.conf \
        /data/conf/nodes.6381.conf \
        /data/conf/nodes.6382.conf \
        /data/conf/nodes.6383.conf \
        /data/conf/nodes.6384.conf ;

    #redisを6台クラスターモードで(クラスターモードの設定はredis.conf)起動。
    # nodes.****.conf はそれぞれ別々のファイルを指定する必要がある。
    redis-server /data/conf/redis.conf --port 6379 --cluster-config-file /data/conf/nodes.6379.conf --daemonize yes ;
    redis-server /data/conf/redis.conf --port 6380 --cluster-config-file /data/conf/nodes.6380.conf --daemonize yes ;
    redis-server /data/conf/redis.conf --port 6381 --cluster-config-file /data/conf/nodes.6381.conf --daemonize yes ;
    redis-server /data/conf/redis.conf --port 6382 --cluster-config-file /data/conf/nodes.6382.conf --daemonize yes ;
    redis-server /data/conf/redis.conf --port 6383 --cluster-config-file /data/conf/nodes.6383.conf --daemonize yes ;
    redis-server /data/conf/redis.conf --port 6384 --cluster-config-file /data/conf/nodes.6384.conf --daemonize yes ;

    REDIS_LOAD_FLG=true;

    #全てのredis-serverの起動が完了するまでループ。
    while $REDIS_LOAD_FLG;
    do
        sleep 1;
        redis-cli -p 6379 info 1> /data/conf/r6379i.log 2> /dev/null;
        if [ -s /data/conf/r6379i.log ]; then
            :
        else
            continue;
        fi
        redis-cli -p 6380 info 1> /data/conf/r6380i.log 2> /dev/null;
        if [ -s /data/conf/r6380i.log ]; then
            :
        else
            continue;
        fi
        redis-cli -p 6381 info 1> /data/conf/r6381i.log 2> /dev/null;
        if [ -s /data/conf/r6381i.log ]; then
            :
        else
            continue;
        fi
        redis-cli -p 6382 info 1> /data/conf/r6382i.log 2> /dev/null;
        if [ -s /data/conf/r6382i.log ]; then
            :
        else
            continue;
        fi
        redis-cli -p 6383 info 1> /data/conf/r6383i.log 2> /dev/null;
        if [ -s /data/conf/r6383i.log ]; then
            :
        else
            continue;
        fi
        redis-cli -p 6384 info 1> /data/conf/r6384i.log 2> /dev/null;
        if [ -s /data/conf/r6384i.log ]; then
            :
        else
            continue;
        fi
        #redis-serverの起動が終わったらクラスター・レプリカの割り当てる。
        #ipを127.0.0.1で割り当てるとphpで不具合が起こるのでpublic ipを指定。
        yes "yes" | redis-cli --cluster create $ipaddr:6379 $ipaddr:6380 $ipaddr:6381 $ipaddr:6382 $ipaddr:6383 $ipaddr:6384 --cluster-replicas 1;
        REDIS_LOAD_FLG=false;
    done
}

if [ ${CLUSTER} -eq 1 ]; then
    cluster
    while true; 
    do
        sleep 1
    done
else
    redis-server /data/conf/redis.conf
fi