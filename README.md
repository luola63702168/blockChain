### 前言
- 主要环境：
    - 语言：go 1.13.5  
    - 数据库：bolt
### 项目构成
   - 实现区块，及区块链基本结构：
   - 矿工工作量难度设置
   - 矿工工作量证明
   - 实现交易稳定可信
   - 实现命令行调用该系统
    
### 命令参考
- createChain --address ADDRESS "create a block Chain"
- printChain			  "print all blocks"
- getBalance --address ADDRESS   "get balance of the address"
- send --from FROM --to TO --amount AMOUNT   "send coin from FROM to TO"

### 项目展示

```shell script
63702@DESKTOP-GONT63B MINGW64 /e/GIT_space/go/blockChain/blockChainV4 (master)
$ go build *.go

63702@DESKTOP-GONT63B MINGW64 /e/GIT_space/go/blockChain/blockChainV4 (master)
$ ./block.exe
invalid input

        createChain --address ADDRESS "create a block Chain"
        printChain                        "print all blocks"
        getBalance --address ADDRESS   "get balance of the address"
        send --from FROM --to TO --amount AMOUNT   "send coin from FROM to TO"


63702@DESKTOP-GONT63B MINGW64 /e/GIT_space/go/blockChain/blockChainV4 (master)
$ ./block.exe  createChain --address "luola"
Begin Mining...
target hash : 010000000000000000000000000000000000000000000000000000000000
found nonce, hash:00000089594836f2355254f70b0f9bdf6fe93e7815e1c26cfb706e71c87327d9,nonce:6149246
Create block successfully

63702@DESKTOP-GONT63B MINGW64 /e/GIT_space/go/blockChain/blockChainV4 (master)
$ ./block.exe  getBalance --address "luola"
The balance of luola is : 12.500000

63702@DESKTOP-GONT63B MINGW64 /e/GIT_space/go/blockChain/blockChainV4 (master)
$ ./block.exe send --from "luola" --to "bob" --amount 2
Begin Mining...
target hash : 010000000000000000000000000000000000000000000000000000000000
found nonce, hash:00000038d21c90d6e8f6df33a584fb5689928bff82acb6ddc8836bf802939dab,nonce:33632683
send successfully!

63702@DESKTOP-GONT63B MINGW64 /e/GIT_space/go/blockChain/blockChainV4 (master)
$ ./block.exe send --from "bob" --to "duck" --amount 1
Begin Mining...
target hash : 010000000000000000000000000000000000000000000000000000000000
found nonce, hash:00000043bf86a8c560f68e9826976fd9620bab9436336ff9d6730abd39baadfc,nonce:8442745
send successfully!

63702@DESKTOP-GONT63B MINGW64 /e/GIT_space/go/blockChain/blockChainV4 (master)
$ ./block.exe getBalance --address "luola"
The balance of luola is : 10.500000

63702@DESKTOP-GONT63B MINGW64 /e/GIT_space/go/blockChain/blockChainV4 (master)
$ ./block.exe getBalance --address "bob"
The balance of bob is : 1.000000

63702@DESKTOP-GONT63B MINGW64 /e/GIT_space/go/blockChain/blockChainV4 (master)
$ ./block.exe printChain
Version1
PreBlockHash:00000038d21c90d6e8f6df33a584fb5689928bff82acb6ddc8836bf802939dab
Hash:00000043bf86a8c560f68e9826976fd9620bab9436336ff9d6730abd39baadfc
TimeStamp:1579691689
Bits:24
Nonce:8442745
IsValid:true
Version1
PreBlockHash:00000089594836f2355254f70b0f9bdf6fe93e7815e1c26cfb706e71c87327d9
Hash:00000038d21c90d6e8f6df33a584fb5689928bff82acb6ddc8836bf802939dab
TimeStamp:1579691619
Bits:24
Nonce:33632683
IsValid:true
Version1
PreBlockHash:
Hash:00000089594836f2355254f70b0f9bdf6fe93e7815e1c26cfb706e71c87327d9
TimeStamp:1579691536
Bits:24
Nonce:6149246
IsValid:true
print over

63702@DESKTOP-GONT63B MINGW64 /e/GIT_space/go/blockChain/blockChainV4 (master)
$
```
![image](https://github.com/luola63702168/blockChain/blob/master/obj_images/index.png)
