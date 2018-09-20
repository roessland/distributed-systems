# Key-Value store with Go net/rpc

## Example run A
```
2018/09/20 20:15:06 KV Get args &{A}
2018/09/20 20:15:06 KV Get reply &{ErrNoKey }
2018/09/20 20:15:06 KV Put args &{A a}
2018/09/20 20:15:06 KV Put args &{A aa}
2018/09/20 20:15:06 KV Put args &{B aaa}
2018/09/20 20:15:06 KV Put args &{B b}
2018/09/20 20:15:06 KV Get args &{B}
2018/09/20 20:15:06 KV Get reply &{OK b}
2018/09/20 20:15:06 KV Put args &{B bb}
2018/09/20 20:15:06 KV Put args &{C bbb}
2018/09/20 20:15:06 KV Put args &{C cc}
2018/09/20 20:15:06 KV Get args &{B}
2018/09/20 20:15:06 KV Get reply &{OK bb}
2018/09/20 20:15:06 KV Put args &{C c}
2018/09/20 20:15:06 KV Put args &{D ccc}
2018/09/20 20:15:06 KV Get args &{B}
2018/09/20 20:15:06 KV Get reply &{OK bb}
2018/09/20 20:15:06 KV Get args &{C}
2018/09/20 20:15:06 KV Get reply &{OK c}
2018/09/20 20:15:06 KV Get args &{D}
2018/09/20 20:15:06 KV Get reply &{OK ccc}
2018/09/20 20:15:06 KV Get args &{D}
2018/09/20 20:15:06 KV Get reply &{OK ccc}
2018/09/20 20:15:06 KV Get args &{A}
2018/09/20 20:15:06 KV Get reply &{OK aa}
Done waiting for waitgroup
```

## Example run B
```
2018/09/20 20:17:34 KV Get args &{A}
2018/09/20 20:17:34 KV Get reply &{ErrNoKey }
2018/09/20 20:17:34 KV Put args &{A a}
2018/09/20 20:17:34 KV Put args &{A aa}
2018/09/20 20:17:34 KV Put args &{D ccc}
2018/09/20 20:17:34 KV Get args &{B}
2018/09/20 20:17:34 KV Get reply &{ErrNoKey }
2018/09/20 20:17:34 KV Put args &{B b}
2018/09/20 20:17:34 KV Put args &{B bb}
2018/09/20 20:17:34 KV Put args &{B aaa}
2018/09/20 20:17:34 KV Get args &{B}
2018/09/20 20:17:34 KV Get reply &{OK aaa}
2018/09/20 20:17:34 KV Put args &{C c}
2018/09/20 20:17:34 KV Put args &{C cc}
2018/09/20 20:17:34 KV Put args &{C bbb}
2018/09/20 20:17:34 KV Get args &{B}
2018/09/20 20:17:34 KV Get reply &{OK aaa}
2018/09/20 20:17:34 KV Get args &{C}
2018/09/20 20:17:34 KV Get reply &{OK bbb}
2018/09/20 20:17:34 KV Get args &{D}
2018/09/20 20:17:34 KV Get reply &{OK ccc}
2018/09/20 20:17:34 KV Get args &{D}
2018/09/20 20:17:34 KV Get reply &{OK ccc}
2018/09/20 20:17:34 KV Get args &{A}
2018/09/20 20:17:34 KV Get reply &{OK aa}
Done waiting for waitgroup

```
