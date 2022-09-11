package main

import (
	"context"
	"go.etcd.io/etcd/client/v3"
	"testing"
	"time"
)

func GetEtcdClient() (*clientv3.Client, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"10.211.55.18:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func TestEtcdKVPut(t *testing.T) {
	// 创建客户端
	client, err := GetEtcdClient()
	if err != nil {
		t.Log(err)
		return
	}
	// 创建读写对象
	kv := clientv3.NewKV(client)

	math, _ := kv.Put(context.TODO(), "/lesson/math", "100")
	t.Log(math)

	english, _ := kv.Put(context.TODO(), "/lesson/english", "90")
	t.Log(english)

	// 获取目录下面所有, 第三个参数with...用来限制，有很多选项
	res, _ := kv.Get(context.TODO(), "/lesson/", clientv3.WithPrefix(), clientv3.WithLimit(2))
	for _, v := range res.Kvs {
		t.Log(string(v.Key))
		t.Log(string(v.Value))
	}
}

func TestEtcdKVDelete(t *testing.T) {
	client, err := GetEtcdClient()
	if err != nil {
		t.Log(err)
		return
	}
	kv := clientv3.NewKV(client)
	// 有这个前缀的批量删除
	_, err = kv.Delete(context.TODO(), "/lesson/",
		clientv3.WithPrefix(),  // 前缀
		clientv3.WithFromKey(), // 按序
	)
	if err != nil {
		return
	}
	//  DO操作
	op := clientv3.OpGet("/lesson/", clientv3.WithLimit(3))
	resp, _ := kv.Do(context.TODO(), op)
	t.Log(resp.Get().Kvs)
}

func TestEtcdKVOp(t *testing.T) {
	client, err := GetEtcdClient()
	if err != nil {
		t.Log(err)
		return
	}
	// 创建读写对象
	kv := clientv3.NewKV(client)

	// OP DO操作
	op := clientv3.OpGet("/lesson/", clientv3.WithLimit(3))
	resp, _ := kv.Do(context.TODO(), op)
	t.Log(resp.Get().Kvs)
}

func TestEtcdLease(t *testing.T) {
	client, err := GetEtcdClient()
	if err != nil {
		return
	}
	kv := clientv3.NewKV(client)
	// 创建租约
	lease := clientv3.NewLease(client)
	leaseR, _ := lease.Grant(context.TODO(), 10)

	// 使用租约put [绑定key]
	putR, _ := kv.Put(context.TODO(), "/corn/lock/job", "100", clientv3.WithLease(leaseR.ID))
	t.Log(putR.Header.Revision)

	// 检查key是否过期
	for {
		getR, err := kv.Get(context.TODO(), "/corn/lock/job")
		if err != nil {
			return
		}
		if getR.Count == 0 {
			t.Log("key过期")
			return
		} else {
			t.Log("key未过期")
			time.Sleep(1 * time.Second)
		}
	}
}

func TestEtcdKeepAliveLease(t *testing.T) {
	client, err := GetEtcdClient()
	if err != nil {
		return
	}
	kv := clientv3.NewKV(client)

	// 创建租约
	lease := clientv3.NewLease(client)
	leaseR, _ := lease.Grant(context.TODO(), 10)

	// 自动续租, 可以超时取消或者可以取消
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 自动续租
	keepChan, _ := lease.KeepAlive(ctx, leaseR.ID)

	// 监听续租情况
	go func() {
		for keepR := range keepChan {
			if keepR == nil {
				t.Log("租约失效")
			} else {
				t.Log("KeepAlive 续约")
			}
		}
	}()
	// 利用租约put
	putR, _ := kv.Put(context.TODO(), "/cron/lock/job", "100", clientv3.WithLease(leaseR.ID))
	t.Log(putR.Header.Revision)

	// 定时查看key是否过期
	go func() {
		for {
			getR, _ := kv.Get(context.TODO(), "/cron/lock/job")
			if getR.Count == 0 {
				t.Log("key过期")
				return
			} else {
				t.Log("key未过期")
				time.Sleep(1 * time.Second)
			}
		}
	}()
	time.Sleep(10 * time.Second)
}

func TestEtcdWatch(t *testing.T) {
	// 创建客户端
	client, err := GetEtcdClient()
	if err != nil {
		return
	}
	// 创建读写对象
	kv := clientv3.NewKV(client)

	// 模拟变化
	go func() {
		for {
			kv.Put(context.TODO(), "/cron/lock/job", "job")
			kv.Delete(context.TODO(), "/cron/lock/job")
			time.Sleep(2 * time.Second)
		}
	}()

	// 获取当前值
	getR, _ := kv.Get(context.TODO(), "/cron/lock/job")
	t.Log(getR)
	t.Log(getR.Header.Revision)
	for _, v := range getR.Kvs {
		t.Log(string(v.Key))
		t.Log(string(v.Value))
	}

	// 监听
	// 当前集群的事务ID revision
	watchStartRevision := getR.Header.Revision + 1
	watcher := clientv3.NewWatcher(client)
	watchChan := watcher.Watch(context.TODO(), "/cron/lock/job", clientv3.WithRev(watchStartRevision))

	for watchChange := range watchChan {
		for _, event := range watchChange.Events {
			switch event.Type {
			case 0:
				t.Log("PUT OP")
			case 1:
				t.Log("DELETE OP")
			}
		}
	}
}

func TestEtcdLock(t *testing.T) {
	client, err := GetEtcdClient()
	if err != nil {
		return
	}
	kv := clientv3.NewKV(client)

	// 以申请租约，创建
	// 申请租约
	lease := clientv3.NewLease(client)
	leaseR, _ := lease.Grant(context.TODO(), 5)

	// 取消租约
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer lease.Revoke(context.TODO(), leaseR.ID)

	// 自动续约
	keepChan, _ := lease.KeepAlive(ctx, leaseR.ID)
	go func() {
		for keepR := range keepChan {
			if keepR == nil {
				t.Log("租约到期")
			} else {
				t.Log("自动续约", leaseR.ID)
			}
		}
	}()
	key := "/cron/lock/job"
	txn := kv.Txn(context.TODO())
	txn.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		// 抢锁成功
		Then(clientv3.OpPut(key, "job", clientv3.WithLease(leaseR.ID))).
		// 抢锁失败
		Else(clientv3.OpGet(key, clientv3.WithLease(leaseR.ID)))

	txnR, _ := txn.Commit()
	if !txnR.Succeeded {
		t.Log("没有抢到锁")
	} else {
		t.Log("抢到锁")
		// 模拟代码
		time.Sleep(time.Second * 5)
		//  释放锁
		cancel()
		t.Log("释放锁")
	}
}
