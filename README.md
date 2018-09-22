# Distributed systems

My notes and code examples related to distributed systems.


## Courses

- [ ] [MIT 6.033 Computer System Engineering](http://web.mit.edu/6.033/www/)
- [ ] [MIT 6.828 Operating System Engineering](https://pdos.csail.mit.edu/6.828/2018/schedule.html)
- [ ] (**in progress**) [MIT 6.824 Distributed Systems 16 week course](http://nil.csail.mit.edu/6.824/2017/)
    - [x] [Lab 1: MapReduce implementation](https://github.com/roessland/distributed-systems/tree/master/mapreduce)
    - [ ] Lab 2: Raft implementation

- [ ] Kubernetes EdX course


## Resources

[Resources on learning system design (Reddit)](https://www.reddit.com/r/cscareerquestions/comments/5u825g/resources_on_learning_system_design_and_data/)

## Books
 - (**in progress, ch1/9**) [Distributed Systems (van Steen, Tanenbaum, 2017)](https://www.distributed-systems.net/)

## Videos/Lectures

- [x] [CS75 (Summer 2012) Lecture 9 Scalability Harvard Web DeveloIpment David Malan](https://www.youtube.com/watch?v=-W9F__D3oY4)
Good summary of standard dns/load balancers/database replication setup.
Focus on typical single points of failure.

- [x] [Distributed Systems in One Lesson by Tim Berglund (50min version)](https://www.youtube.com/watch?v=Y6Ev8GIlbxc) -- great introduction to what we give up to achieve horizontal scaling.

- [x] [Designing for Understandability: The Raft Consensus Algorithm](https://www.youtube.com/watch?v=vYp4LYbnnW8) - includes demo of Raft visualization usage

## Papers
- [x] [MapReduce (2004)](http://nil.csail.mit.edu/6.824/2018/papers/mapreduce.pdf)

## Blogs

- [High Scalability](http://highscalability.com/)
- [DBMS Musings](http://dbmsmusings.blogspot.com/2018/09/newsql-database-systems-are-failing-to.html)
    - [x] [NewSQL database systems are failing to guarantee consistency, and I blame Spanner](http://dbmsmusings.blogspot.com/2018/09/newsql-database-systems-are-failing-to.html)

## Keywords and notes

- [x] Scalability
- [x] Vertical scaling (Add more CPU, RAM, disks, etc. Limited, a CPU can only go so fast.)
- [x] Horizontal scaling (Add more computers, more app instances, more database instances, more message queue brokers. Adds complexity and takes away good things, see CAP theorem. Commodity hardware.)
- [x] DNS (Domain Name Server, returns IP)
- [ ] Load balancing
    - [ ] Heart beats and health checks
    - [ ] Random-ish distribution using hashing
    - [ ] Reverse proxy load balancer
        - [ ] HAProxy, ELB, LVS
    - [ ] Hardware load balancer (enterprise hardware)
        - [ ] Barracuda, F5, Cisco, Citrix (expensive stuff!)
        - active/active-mode enables redundancy
    - [ ] DNS load balancer (BIND)
    - [ ] Round-robin (every nth request goes to nth server)
    - Can add and read cookies, using e.g. a random number to send user to same
      server as the last request.
    - [ ] How to prevent load balancer failure?
    - security -- prevent load balancer bypass, firewalls
    - load balancer (before nginx) can do SSL termination
- [ ] Security
    - [ ] Firewalls -- least privilege principle
- [ ] Caching
    - memcached
    - redis
        - [ ] distributed redis (sharding, replication, scaling)
    - cache eviction (TTL, LRU, FIFO)
- [ ] HTTP Sessions (need a shared session store or JWTs)
    - [ ] Session cookies (contains session id)
    - [ ] Secure cookies (encrypted)
    - [ ] JWT (contains signed proof of access to a resource, hard to revoke)
- [ ] Shared storage (HDFS, Gluster, Ceph, NFS)
    - [ ] RAID (use it)
- [ ] Database scaling
    - [ ] Database replication
        - Can use a load balancer in front of all replicas
    - [ ] Database sharding
    - [ ] Partitioning
- [ ] Faults
    - [ ] Fault injection
        - [ ] Netflix stack for fault injection
- [ ] Failures
- [ ] Fault tolerance
    - [ ] Timeouts
    - [ ] Retries
    - [ ] Circuit breakers
    - [ ] Fail fast
    - [ ] Fallbacks
- [ ] Rate limiting
    - [ ] Shared rate limit (e.g. cluster of services accessing same third party API)
- [ ] Monitoring
    - [ ] Prometheus -- how does it scale? What is max retention?
    - [ ] InfluxDB -- how does it scale? What is max retention?
- [ ] Logging
- [ ] [Mesos](http://mesos.apache.org/)
- [ ] Automated deployment
    - [ ] Ansible
- [ ] Containers
    - [ ] Docker
- [ ] Container orchestration
    - [ ] Kubernetes
- [ ] Distributed key value stores
    - [ ] Distributed agreement protocols
        - [ ] Raft - Replicated And Fault Tolerant (!)
            - [x] https://raft.github.io/ Good visualization (RaftScope)
            - [x] [Students guide to raft (TheSquarePlanet)](https://thesquareplanet.com/blog/students-guide-to-raft/)
            - [x] http://thesecretlivesofdata.com/raft/
            - [ ] [Lecture from Raft User Study](https://www.youtube.com/watch?v=YbZ3zDzDnrw)
        - [ ] Paxos
    - [ ] Quorum
    - [ ] Zookeeper
- [ ] Distributed databases
    - [ ] Cassandra
    - [ ] Sharded/replicated PostgreSQL
    - [ ] Replication
        - [ ] A: For reads only (agreement not a problem, doesn't scale forever)
        - [ ] B: For writes and reads (scales, but there are cons)
    - [ ] Consistent hashing
    - [ ] Network partitioning (what to do for reads? for writes?)
- [ ] Microservices pattern
- [ ] Messages
    - [ ] Message queues
        - [ ] Kafka
    - [ ] Pubsub
    - [ ] Exactly-once delivery (not guaranteed)
    - [ ] At-least-once delivery
    - [ ] Idempotency (as a consequence of at-least-once delivery this is needed)
    - [ ] Encoding
        - [ ] Protobuf
        - [ ] JSON
        - [ ] gRPC
- [ ] P2P
    - [ ] BitCoin
    - [ ] Kazaa
    - [ ] BitTorrent
- [ ] Java (seems to be popular for distributed software, e.g. Apache projects)
- [ ] Distributed parallel computation (HPC)
    - [ ] MPI (when a single computation doesn't fit in the memory of a single computer)
    - [ ] GPU, CUDA, OpenCL, heterogenous computation. (A bit offtopic.)
- [ ] Testing
- [ ] Managed cloud solutions to distributed problems
    - [ ] Azure Service Bus
    - [ ] Vendor lock-in vs do-it-yourself (pros, cons?)
    - [ ] S3 / other blob storage
- [ ] DIY Cloud
    - [ ] [CNCF Cloud Native Interactive Landscape](https://landscape.cncf.io/) (Open source solutions to containers, orchestration, service mesh, discovery, distributed database, ci/cd, monitoring, tracing, logging, networking, messages)
- [ ] Distributed high-latency systems
    - [ ] Georeplication
    - [ ] Datacenter load balancing on DNS level -- mywebsite.com returns IP to
      load balancer in one of several datacenters.
    - [ ] Multi-datacenter redundancy -- tolerate an entire datacenter falling
          off the grid
    - [ ] AWS availability zones (zone with multiple datacenters)
- [ ] Case studies. Architecture, scaling. 
    - [ ] Reddit 
    - [ ] Twitter
    - [ ] YouTube
    - [x] MapReduce
