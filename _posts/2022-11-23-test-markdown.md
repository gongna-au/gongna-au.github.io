---
layout: post
title: 文件系统？
subtitle:
tags: [文件系统]
---

# 文件系统

## 1.文件系统的组成？

文件系统 = **操作系统的子系统**=**1.负责把用户文件放到硬盘**（放到硬盘的数据不会丢失，所以可以持久化的保存数据。）**2.管理放到磁盘上面的所有文件**（管理方式『组织文件的方式』不同，所以文件系统也就不同）

文件系统的基本数据单位？
文件系统的基本数据单位是：“文件”，Linux 最经典的⼀句话是：「⼀切皆⽂件」，不仅普通的⽂件和⽬录，就连块设备、管道、socket 等，也都是统⼀交给⽂件系统管理的。

文件系统的基本操作单位？
⽂件系统的基本操作单位是数据块。

Linux 文件系统管理（磁盘中）文件的方式？
Linux 文件系统为每个文件分配两个数据结构：**1.索引节点『IndexNode』** ,**2.目录项『Dicrectory entry』**。索引节点记录文件的元信息，目录项目记录文件的层次结构

- 索引节点『IndexNode』: IndexNode 用来记录文件的元信息，比如有编号，文件大小，访问权限，创建时间，修改时间，数据在磁盘上面的位置，索引节点是文件的唯一标志符，索引节点同样被存储在磁盘当中（占据磁盘空间）。
- 目录项『Dicrectory entry』: 用来记录文件的名字，索引节点指针，以及与其他目录项的层级关联关系。多个目录项目关联起来就是目录结构，但他和索引节点不同的是，它缓存在内存中间，由内核维护的一个数据结构。

```go

// 超级块区 存储文件系统的详细信息。块个数，块大小，空闲块 ，在文件系统挂载的时候进入内存
type SuperBlock struct{
    // 块个数
    BlockNum int
    // 块大小
    BlockSize int
    // 空闲块
    SpareBlock int
}

// 数据块 一个索引节点指向一个数据块
type DataBlock struct{
    // 数据块的位置
    StartLocation int
    EndLocation int
}

// 当文件被访问的时候加载进内存
type IndexNodeBlock struct{
    ID int
    // 创建时间
    CreateTime time.Time
    // 修改时间
    UpdateTime time.Time
    // 文件大小
    FileSize int
    // 访问权限
    Priviledge string
    DataBlock  *DataBlock

}

type Dicrectory struct{
    Nodes []*IndexNodeBlock
    // 一个目录下面可以有多个节点，所以这里是切片
    Son   []*Dicrectory
    SonNum int
}

```

- **一个索引节点存储着一个数据块，一个目录项又存储着多个索引节点。**
- 索引节点和数据块都在磁盘（硬盘）中，但是目录项是存储在内存中的。
- 为了查找数据更加的快速，有时由会把索引节点也从硬盘中间加再到内存中。
- 磁盘在格式化的时候，会分成三个存储区域：1.超级块。2.索引节点区 3.数据区
- 超级块：文件系统挂载进入内存，索引节点块：当文件被访问的时候进入内存。

## 2.虚拟系统？

文件系统由很多，但是操作系统希望可以有一个统一的文件接口，不管什么样的文件系统都可以被操作系统使用，于是，就有了虚拟文件系统 VFS。

```go
type VitualFileSystem interface{
    Schedule()
}
// 磁盘文件系统
type DiskFileSystem struct{
    Buffer *Buffer
}

func (c *ConcreteFileSystem) Schedule(){

}

// 内存文件系统
type InMemoryFileSystem struct{
    Buffer *Buffer
}

func (i *InMemoryFileSystem ) Schedule(){

}

// 网络文件系统
type NetworkFileSystem struct{
    Buffer *Buffer
}

func (n *NetworkFileSystem) Schedule(){

}

type Buffer struct{

}

```

Linux ⽀持的⽂件系统也不少，根据存储位置的不同，可以把⽂件系统分为三类：

- 1.磁盘的⽂件系统，它是直接把数据存储在磁盘中，⽐如 Ext 2/3/4、XFS 等都是这类⽂件系统。
- 2.内存的⽂件系统，这类⽂件系统的数据不是存储在硬盘的，⽽是占⽤内存空间，我们经常⽤到的 /proc 和 /sys ⽂件系统都属于这⼀类，读写这类⽂件，实际上是读写内核中相关的数据。
- 3.⽹络的⽂件系统，⽤来访问其他计算机主机数据的⽂件系统，⽐如 NFS、SMB 等等。
  > ⽂件系统⾸先要先挂载到某个⽬录才可以正常使⽤，⽐如 Linux 系统在启动时，会把⽂件系统挂载到根⽬录。
  > Linux 在启动的时候，会把文件系统挂载到根目录

## 3.文件的使用？

```go
type FILE struct{

}

func OpenFile(filename string)*FILE{

}

// 写⽂件
func WriteFile(f *FILE，path string){
    //
}

// 读⽂件
func ReadFile(f *FILE){
    //
}

// 关闭⽂件
func CloseFile(f *FILE){
    //
}


```

1.⾸先⽤ open 系统调⽤打开⽂件， open 的参数中包含⽂件的路径名和⽂件名。2.使⽤ write 写数据，其中 write 使⽤ open 所返回的⽂件描述符，并不使⽤⽂件名作为参数。3.使⽤完⽂件后，要⽤ close 系统调⽤关闭⽂件，避免资源的泄露

操作系统如何跟踪某个进程打开的文件？

我们打开了⼀个⽂件后，操作系统会跟踪进程打开的所有⽂件，所谓的跟踪呢，就是操作系统为每个进程**维护⼀个打开⽂件表**，⽂件表⾥的每⼀项代表「⽂件描述符」，所以说⽂件描述符是打开⽂件的标识。

```go
type Process struct{
    FILEsOPenedTable *FILEsOPenedTable
}

type FILEsOPenedTable struct{
    FILEsOPened []*FILE
}

type FILE struct{
    // 这个文件指针每被使用一次，就加1，标志这个文件被多少个进程打开，只有当这个Counter==0的时候，才能执行关闭文件，释放资源的函数
    Counter int
    // 文件的磁盘位置
    Location int
    // 这个文件允许的操作
    Priviledge int
}
// 写⽂件
func (f *FILE) WriteFile(path string){
}

// 读⽂件
func (f *FILE)  ReadFile(){

}

// 关闭⽂件
func (f *FILE)  CloseFile(){
    if f.Counter ==0{
        //
    }
}

```

操作系统在打开⽂件表中维护着打开⽂件的状态和信息：

- ⽂件指针：系统跟踪上次读写位置作为当前⽂件位置指针，这种指针对打开⽂件的某个进程来说是唯⼀的；
- ⽂件打开计数器：⽂件关闭时，操作系统必须重⽤其打开⽂件表条⽬，否则表内空间不够⽤。因为多个进程可能打开同⼀个⽂件，所以系统在删除打开⽂件条⽬之前，必须等待最后⼀个进程关闭⽂件，该计数器跟踪打开和关闭的数量，当该计数为 0 时，系统关闭⽂件，删除该条⽬；
- ⽂件磁盘位置：绝⼤多数⽂件操作都要求系统修改⽂件数据，该信息保存在内存中，以免每个操作都从磁盘中读取；
- 访问权限：每个进程打开⽂件都需要有⼀个访问模式（创建、只读、读写、添加等），该信息保存在进程的打开⽂件表中，以便操作系统能允许或拒绝之后的 I/O 请求；

操作系统的视⻆是如何把⽂件数据和磁盘块对应起来?
⽤户和操作系统对⽂件的读写操作是有差异的，⽤户习惯以字节的⽅式读写⽂件，⽽操作系统则是以数据块来读写⽂件，那屏蔽掉这种差异的⼯作就是⽂件系统了。

我们来分别看⼀下，读⽂件和写⽂件的过程?
读取文件：“当用户进程读取一个字节大小的数据的时候，文件系统则需要获取到字节所对应的数据块，在返回数据块对应的数据部分”
写文件：用户进程把写一个字节大小的数据写进文件的时候，文件找到需要写入数据的数据块的位置，然后修改数据块的数据部分，最后把数据写回磁盘。

## 4.文件的存储？

⽂件的数据是要存储在硬盘上⾯的，数据在磁盘上的存放⽅式，就像程序在内存中存放的⽅式那样，有以下两种：

- 1.连续空间存放⽅式
- 2.⾮连续空间存放⽅式
  其中，⾮连续空间存放⽅式⼜可以分为「链表⽅式」和「索引⽅式」。

### 『连续空间存储』

> 前提是知道文件的大小（文件头通过起始的位置和长度）注意，此处说的⽂件头，就类似于 Linux 的 inode。

```go
type FileHead struct{
    StartIndex int
    Legth int
}
```

> 连续空间存放⽅式顾名思义，⽂件存放在磁盘「连续的」物理空间中。这种模式下，⽂件的数据都是紧密相连，读写效率很⾼，因为⼀次磁盘寻道就可以读出整个⽂件。

使⽤连续存放的⽅式有⼀个前提，必须先知道⼀个⽂件的⼤⼩，这样⽂件系统才会根据⽂件
的⼤⼩在磁盘上找到⼀块连续的空间分配给⽂件。

所以，⽂件头⾥需要指定「起始块的位置」和「⻓度」，有了这两个信息就可以很好的表示
⽂件存放⽅式是⼀块连续的磁盘空间。

> 缺陷 1：「磁盘空间碎⽚」
> 连续空间存放的⽅式虽然读写效率⾼，但是有「磁盘空间碎⽚」和「⽂件⻓度不易扩展」的缺陷。
> 如下图，如果⽂件 B 被删除，磁盘上就留下⼀块空缺，这时，如果新来的⽂件⼩于其中的⼀个空缺，我们就可以将其放在相应空缺⾥。但如果该⽂件的⼤⼩⼤于所有的空缺，但却⼩于空缺⼤⼩之和，则虽然磁盘上有⾜够的空缺，但该⽂件还是不能存放。当然了，我们可以通过将现有⽂件进⾏挪动来腾出空间以容纳新的⽂件，但是这个在磁盘挪动⽂件是⾮常耗时，所以这种⽅式不太现实。
> 「⽂件⻓度不易扩展」
> 另外⼀个缺陷是⽂件⻓度扩展不⽅便，例如上图中的⽂件 A 要想扩⼤⼀下，需要更多的磁盘空间，唯⼀的办法就只能是挪动的⽅式，前⾯也说了，这种⽅式效率是⾮常低的。

```go
type Disk struct{
   Data []int
}

func (d *Disk)DeleteFile(start int ,end int){
    for i:=start ;i<=end ;i++{
        d.Data[i]=0
    }
}
```

### 『非连续空间存储』

⾮连续空间存放⽅式分为「链表⽅式」和「索引⽅式」。

> 「链表⽅式」-隐式链接
> ⽂件要以「隐式链表」的⽅式存放的话，实现的⽅式是⽂件头要包含「第⼀块」和「最后⼀块」的位置，并且每个数据块⾥⾯留出⼀个指针空间，⽤来存放下⼀个数据块的位置，这样⼀个数据块连着⼀个数据块，从链头开是就可以顺着指针找到所有的数据块，所以存放的⽅式可以是不连续的。隐式链表的存放⽅式的缺点在于⽆法直接访问数据块，只能通过指针顺序访问⽂件，以及数据块指针消耗了⼀定的存储空间。隐式链接分配的稳定性较差，系统在运⾏过程中由于软件或者硬件错误导致链表中的指针丢失或损坏，会导致⽂件数据的丢失。

Tips:隐式链接是每个数据块存放下一个数据块的位置。

> 「链表⽅式」-显式链接
> 对于显式链接的⼯作⽅式，我们举个例⼦，⽂件 A 依次使⽤了磁盘块 4、7、2、10 和 12 ，⽂件 B 依次使⽤了磁盘块 6、3、11 和 14 。利⽤下图中的表，可以从第 4 块开始，顺着链⾛到最后，找到⽂件 A 的全部磁盘块。同样，从第 6 块开始，顺着链⾛到最后，也能够找出⽂件 B 的全部磁盘块最后，这两个链都以⼀个不属于有效磁盘编号的特殊标记（如 -1 ）结束。内存中的这样⼀个表格称为⽂件分配表（File Allocation Table ， FAT）。

> 「索引⽅式」
> 链表的⽅式解决了连续分配的磁盘碎⽚和⽂件动态扩展的问题，但是不能有效⽀持直接访问 FAT 除外），索引的⽅式可以解决这个问题。索引的实现是为每个⽂件创建⼀个「索引数据块」，⾥⾯存放的是指向⽂件数据块的指针列表，说⽩了就像书的⽬录⼀样，要找哪个章节的内容，看⽬录查就可以。⽂件头需要包含指向「索引数据块」的指针切片，这样就可以通过⽂件头知道所有数据块的位置，再通过具体的数据块指针，找到具体的数据块。

```go
type FILE struct{
    Head FileHead
}

type FileHead struct{
    // locations 存储的是 DataBlock 在Disk中间的下标
    Locations []int
}

type DataBlock struct{
    IDataContent string
}
// 模拟物理磁盘中按照数据块存储数据
type Disk struct{
    DataBlocks []DataBlock
}
```

Tips:因为从文件头可以知道所有的数据块在 Disk 的位置，所以可以直接 Disk.DataBlocks[i]的方式直接获取到文件

- ⽂件的创建、增⼤、缩⼩很⽅便
- 不会有碎⽚的问题
- ⽀持顺序读写和随机读写；
- 由于索引数据也是存放在磁盘块的，如果⽂件很⼩，只需⼀块就可以存放的下，但还是需要额外分配⼀块来存放索引数据，所以缺陷之⼀就是存储索引带来的开销。

```go
func (f *FILE) AddFile(index int){
    f.Head.Locations=append(f.Head.Locations,index)
}
```

如何通过组合的方式处理大文件的存放？

如果⽂件很⼤，⼤到⼀个索引数据块放不下索引信息，这时⼜要如何处理⼤⽂件的存放呢？我们可以通过组合的⽅式，来处理⼤⽂件的存储。

> 「链式索引』
> 先来看看链表 + 索引的组合，这种组合称为「链式索引块」，它的实现⽅式是在索引数据块留出⼀个存放下⼀个索引数据块的指针。于是当⼀个索引数据块的索引信息⽤完了，就可以通过指针的⽅式，找到下⼀个索引数据块的信息。那这种⽅式也会出现前⾯提到的链表⽅式的问题，万⼀某个指针损坏了，后⾯的数据也就会⽆法读取了。

```go
type FILE struct{
    Head FileHead
}

type FileHead struct{
    // locations 存储的是 DataBlock 在Disk中间的下标
    Locations []Location
}

type Location struct{
    Index int
    Next  int
}

type DataBlock struct{
    IDataContent string
}
// 模拟物理磁盘中按照数据块存储数据
type Disk struct{
    DataBlocks []DataBlock
}

```

> 索引 + 索引的⽅式，这种组合称为「多级索引块』

『Unix ⽂件的实现⽅式』
文件实现方式对比：

顺序分配：

- 访问磁盘一次。顺序分配存取速度块，当文件式定长的时候，可以根据文件的起始地址和长度进行随机访问。连续的物理存储空间。所有的数据块可以一次性返回读出

链表分配：

- 访问磁盘 n 次。因为每次访问一个数据块，所以返回一个数据块，并且每次只能知道下一个数据块的位置，其他的数据块的位置是不知道的。

索引分配：

- M 级需要访问磁盘 m+1 次，因为每次的每一个级可以读出所有的数据。

那早期 Unix ⽂件系统是组合了前⾯的⽂件存放⽅式的优点它是根据⽂件的⼤⼩，存放的⽅式会有所变化：

```go

type FILEIndexNode struct {
	NodesLevel []IndexNodes
}

type IndexNodes struct {
	Indexs []int
	Next   *IndexNodes
}

func Test() {
	var nodes FILEIndexNode
	nodes = FILEIndexNode{}
	nodes.NodesLevel = make([]IndexNodes, 4)
	nodes.NodesLevel[0] = IndexNodes{
		Indexs: []int{1, 2, 3, 4, 5, 6, 7},
		Next:nil,
	}
	nodes.NodesLevel[1] = IndexNodes{
		Next: &IndexNodes{
			Indexs: []int{
				1, 2, 3, 4, 5, 6, 7, 8,
			},
		},
	}

	nodes.NodesLevel[2] = IndexNodes{
		Next: &IndexNodes{
			Next: &IndexNodes{
				Indexs: []int{
					1, 2, 3, 4, 5, 6,
				},
			},
		},
	}

	nodes.NodesLevel[3] = IndexNodes{
		Next: &IndexNodes{
			Next: &IndexNodes{
				Next: &IndexNodes{
					Indexs: []int{
						1, 2, 3, 4, 5, 6,
					},
				},
			},
		},
	}
}

```

- 如果存放⽂件所需的数据块⼩于 10 块，则采⽤直接查找的⽅式；
- 如果存放⽂件所需的数据块超过 10 块，则采⽤⼀级间接索引⽅式；
- 如果前⾯两种⽅式都不够存放⼤⽂件，则采⽤⼆级间接索引⽅式；
- 如果⼆级间接索引也不够存放⼤⽂件，这采⽤三级间接索引⽅式；

所以，这种⽅式能很灵活地⽀持⼩⽂件和⼤⽂件的存放：对于⼩⽂件使⽤直接查找的⽅式可减少索引数据块的开销；对于⼤⽂件则以多级索引的⽅式来⽀持，所以⼤⽂件在访问数据块时需要⼤量查询；这个⽅案就⽤在了 Linux Ext 2/3 ⽂件系统⾥，虽然解决⼤⽂件的存储，但是对于⼤⽂件的访
问，需要⼤量的查询，效率⽐较低。

## 5.文件的空闲空间的管理？

> 简单来说，空闲空间的管理，就是决定如何把一个新的数据存放在硬盘的某个位置。

前⾯说到的⽂件的存储是针对已经被占⽤的数据块组织和管理，接下来的问题是，如果我要保存⼀个数据块，我应该放在硬盘上的哪个位置呢？难道需要将所有的块扫描⼀遍，找个空的地⽅随便放吗？那这种⽅式效率就太低了，所以针对磁盘的空闲空间也是要引⼊管理的机制，接下来介绍⼏种常⻅的⽅法：

- 空闲表法
- 空闲链表法
- 位图法

### 『空闲表』

> 空闲表法
> 空闲表法就是为所有空闲空间建⽴⼀张表，表内容包括空闲区的第⼀个块号和该空闲区的块个数，注意，这个⽅式是连续分配的。如下代码所展示：

```go
type DisengagedBlockTable struct{
    Blocks []ItemBlock
}

type ItemBlock struct{
    StartBlockIndex int
    Legth int
}
```

当请求分配磁盘空间时，系统依次扫描空闲表⾥的内容，直到找到⼀个合适的空闲区域为⽌。当⽤户撤销⼀个⽂件时，系统回收⽂件空间。这时，也需顺序扫描空闲表，寻找⼀个空闲表条⽬并将释放空间的第⼀个物理块号及它占⽤的块数填到这个条⽬中。这种⽅法仅当有少量的空闲区时才有较好的效果。因为，如果存储空间中有着⼤量的⼩的空闲区，则空闲表变得很⼤，这样查询效率会很低。另外，这种分配技术适⽤于建⽴连续⽂件。

### 『空闲链』

> 空闲链表法
> 空闲链表法我们也可以使⽤「链表」的⽅式来管理空闲空间，每⼀个空闲块⾥有⼀个指针指向下⼀个空闲块，这样也能很⽅便的找到空闲块并管理起来。如下代码所示：当创建⽂件需要⼀块或⼏块时，就从链头上依次取下⼀块或⼏块。反之,当回收空间时，把这些空闲块依次接到链头上。

```go

type DisengagedBlockTable struct {
	List *list.List
}

type Item struct {
	Index int
}

func NewDisengagedBlockTable() *DisengagedBlockTable {
	return &DisengagedBlockTable{
		List: list.New(),
	}
}

func (d *DisengagedBlockTable) Add(item *Item) {
	d.List.PushFront(item)
}

func (d *DisengagedBlockTable) Pop(num int) []*Item {
	if d.List.Len() >num {
		result := []*Item{}
		for i:=0;i<num;i++{
			head := d.List.Front()
			if v,ok:=head.Value.(*Item);ok{
				result=append(result,v)
			}
		}
		return result
	}
	fmt.Println("There are not enough blocks")
	return []*Item{}
}
```

这种技术只要在主存中保存⼀个指针，令它指向第⼀个空闲块。其特点是简单，但不能随机访问，⼯作效率低，因为每当在链上增加或移动空闲块时需要做很多 I/O 操作，同时数据块的指针消耗了⼀定的存储空间。空闲表法和空闲链表法都不适合⽤于⼤型⽂件系统，因为这会使空闲表或空闲链表太⼤。

### 『位图法』

> 位图法

位图是利⽤⼆进制的⼀位来表示磁盘中⼀个盘块的使⽤情况，磁盘上所有的盘块都有⼀个⼆进制位与之对应。当值为 0 时，表示对应的盘块空闲，值为 1 时，表示对应的盘块已分配。它形式如下：

```shell
1111110011111110001110110111111100111 ...
```

在 Linux ⽂件系统就采⽤了位图的⽅式来管理空闲空间，不仅⽤于数据空闲块的管理，还⽤于 inode 空闲块的管理，因为 inode 也是存储在磁盘的，⾃然也要有对其管理。

## 6.文件系统的结构？

前⾯提到 Linux 是⽤位图的⽅式管理空闲空间，⽤户在创建⼀个新⽂件时，Linux 内核会通过 inode 的位图找到空闲可⽤的 inode，并进⾏分配。要存储数据时，会通过块的位图找到空闲的块，并分配，但仔细计算⼀下还是有问题的。

1.数据块的位图是放在磁盘块⾥的。 2.如果数据块放在一个磁盘块里面，一个磁盘块 4K 那么最大只能存储 4* 1024* 8= 2^15 个空闲块 由于 1 个数据块是 4K ⼤⼩，那么最⼤可以表示的空间为 2^15 _ 4 _ 1024 = 2^27 个 byte，也就是 128M。
也就是说按照上⾯的结构，如果采⽤「⼀个块的位图 + ⼀系列的块」，外加「⼀个块的 inode 的位图 + ⼀系列的 inode 的结构」能表示的最⼤空间也就 128M，这太少了现在很多⽂件都⽐这个⼤。
在 Linux ⽂件系统，把这个结构称为⼀个块组，那么有 N 多的块组，就能够表示 N ⼤的⽂件。下图给出了 Linux Ext2 整个⽂件系统的结构和块组的内容，⽂件系统都由⼤量块组组成，在硬盘上相继排布：

引导块--块组 1--块组 1--块组 2 --块组 n

> Innode 是索引节点
> 每一个块组的构成：
> 超级块--块组描述符--数据位图--Innode 位图--Innode 列表-- 数据块
> 最前⾯的第⼀个块是引导块，在系统启动时⽤于启⽤引导，接着后⾯就是⼀个⼀个连续的块组了，块组的内容如下：

- 超级块，包含的是⽂件系统的重要信息，⽐如 inode 总个数、块总个数、每个块组的 inode 个数、每个块组的块个数等等。
- 块组描述符，包含⽂件系统中各个块组的状态，⽐如块组中空闲块和 inode 的数⽬等，每个块组都包含了⽂件系统中「所有块组的组描述符信息」
- 数据位图和 inode 位图， ⽤于表示对应的数据块或 inode 是空闲的，还是被使⽤中。
- inode 列表，包含了块组中所有的 inode，inode ⽤于保存⽂件系统中与各个⽂件和⽬录相关的所有元数据。
- 数据块，包含⽂件的有⽤数据。

可以会发现每个块组⾥有很多重复的信息，⽐如超级块和块组描述符表，这两个都是全局信息，⽽且⾮常的重要，这么做是有两个原因：如果系统崩溃破坏了超级块或块组描述符，有关⽂件系统结构和内容的所有信息都会丢失。如果有冗余的副本，该信息是可能恢复的。通过使⽂件和管理数据尽可能近，减少了磁头寻道和旋转，这可以提⾼⽂件系统的性能。
不过，Ext2 的后续版本采⽤了稀疏技术。该做法是，超级块和块组描述符表不再存储到⽂件系统的每个块组中，⽽是只写⼊到块组 0、块组 1 和其他 ID 可以表示为 3、 5、7 的幂的块组中。

## 7.目录的存储？

Tips:目录也是文件。（存储在磁盘）InNode(索引节点)块里面的内容是指向具体的数据块。（存储在磁盘）目录块里面存储的内容是一项一项的文件信息。

在前⾯，我们知道了⼀个普通⽂件是如何存储的，但还有⼀个特殊的⽂件，经常⽤到的⽬录，它是如何保存的呢？基于 Linux ⼀切皆⽂件的设计思想，⽬录其实也是个⽂件，甚⾄可以通过 vim 打开它，它也有 inode，inode ⾥⾯也是指向⼀些块。和普通⽂件不同的是，普通⽂件的块⾥⾯保存的是⽂件数据，⽽⽬录⽂件的块⾥⾯保存的是⽬录⾥⾯⼀项⼀项的⽂件信息。

### 『列表』

在⽬录⽂件的块中，最简单的保存格式就是列表，就是⼀项⼀项地将⽬录下的⽂件信息（如⽂件名、⽂件 inode、⽂件类型等）列在表⾥。列表中每⼀项就代表该⽬录下的⽂件的⽂件名和对应的 inode，通过这个 inode，就可以找到真正的⽂件。

所以(InNode)索引节点 和目录项，以及列表，他们三者以及真正存储数据的数据块)之间的关系是：

```go
// 索引节点
type IndexNode struct{
    DataBlock *DataBlock
}

// 数据块
type DataBlock struct{
    Content string
}

// 目录项
type DirectoryEntry struct{
    List  *List

}

type List struct{
   IndexNodes []*IndexNode
}


```

### 『哈希表』

如果⼀个⽬录有超级多的⽂件，我们要想在这个⽬录下找⽂件，按照列表⼀项⼀项的找，效率就不⾼了。
于是，保存⽬录的格式改成哈希表，对⽂件名进⾏哈希计算，把哈希值保存起来，如果我们要查找⼀个⽬录下⾯的⽂件名，可以通过名称取哈希。如果哈希能够匹配上，就说明这个⽂件的信息在相应的块⾥⾯。
Linux 系统的 ext ⽂件系统就是采⽤了哈希表，来保存⽬录的内容，这种⽅法的优点是查找⾮常迅速，插⼊和删除也较简单，不过需要⼀些预备措施来避免哈希冲突。
⽬录查询是通过在磁盘上反复搜索完成，需要不断地进⾏ I/O 操作，开销较⼤。所以，为了
减少 I/O 操作，把当前使⽤的⽂件⽬录缓存在内存，以后要使⽤该⽂件时只要在内存中操作，从⽽降低了磁盘操作次数，提⾼了⽂件系统的访问速度。

## 8.软链接个硬链接

有时候我们希望给某个⽂件取个别名，那么在 Linux 中可以通过硬链接（Hard Link） 和软链接（Symbolic Link） 的⽅式来实现，它们都是⽐较特殊的⽂件，但是实现⽅式也是不相同的。

硬链接是多个⽬录项中的「索引节点」指向⼀个⽂件，也就是指向同⼀个 inode，但是 inode 是不可能跨越⽂件系统的，每个⽂件系统都有各⾃的 inode 数据结构和列表，所以硬链接是
不可⽤于跨⽂件系统的。由于多个⽬录项都是指向⼀个 inode，那么**只有删除⽂件的所有硬链接以及源⽂件时，系统才会彻底删除该⽂件**。

软链接相当于重新创建⼀个⽂件，这个⽂件有独⽴的 inode，但是这个⽂件的内容是另外⼀个⽂件的路径，所以访问软链接的时候，实际上相当于访问到了另外⼀个⽂件，所以软链接是可以跨⽂件系统的，甚⾄⽬标⽂件被删除了，链接⽂件还是在的，只不过指向的⽂件找不到了⽽已。

## 9.文件 I/O

### 『缓冲和非缓冲』

⽂件操作的标准库是可以实现数据的缓存，那么根据**是否利⽤标准库缓冲**，可以把⽂件 I/O 分为缓冲 I/O 和⾮缓冲 I/O：

- 缓冲 I/O:利⽤的是标准库的缓存实现⽂件的加速访问，⽽标准库再通过系统调⽤访问⽂件。
- ⾮缓冲 I/O，直接通过系统调⽤访问⽂件，不经过标准库缓存.

这⾥所说的「缓冲」特指标准库内部实现的缓冲。⽐⽅说，很多程序遇到换⾏时才真正输出，⽽换⾏前的内容，其实就是被标准库暂时缓存了起来，这样做的⽬的是，减少系统调⽤的次数，毕竟系统调⽤是有 CPU 上下⽂切换的开销的。

### 『直接和非直接』

我们都知道磁盘 I/O 是⾮常慢的，所以 Linux 内核为了减少磁盘 I/O 次数，在系统调⽤后，会把⽤户数据拷⻉到内核中缓存起来，这个内核缓存空间也就是「⻚缓存」，只有当缓存满⾜某些条件的时候，才发起磁盘 I/O 的请求。

根据是「**是否否利⽤操作系统的缓存**」，可以把⽂件 I/O 分为直接 I/O 与⾮直接 I/O：

- 直接 I/O：不会发⽣内核缓存和⽤户程序之间数据复制，⽽是直接经过⽂件系统访问磁盘。
- ⾮直接 I/O，读操作时，数据从**内核缓存**中拷⻉给⽤户程序，写操作时，数据从⽤户程序拷⻉给**内核缓存**，再由内核决定什么时候写⼊数据到磁盘。

如果在使⽤⽂件操作类的系统调⽤函数时，指定了 O_DIRECT 标志，则表示使⽤直接 I/O。。如果没有设置过，默认使⽤的是⾮直接 I/O。如果⽤了⾮直接 I/O 进⾏写数据操作，内核什么情况下才会把缓存数据写⼊到磁盘？

以下⼏种场景会触发内核缓存的数据写⼊磁盘?

- 在调⽤ write 的最后，当发现内核缓存的数据太多的时候，内核会把数据写到磁盘上；
- ⽤户主动调⽤ sync ，内核缓存会刷到磁盘上；
- 当内存⼗分紧张，⽆法再分配⻚⾯时，也会把内核缓存的数据刷到磁盘上；
- 内核缓存的数据的缓存时间超过某个时间时，也会把数据刷到磁盘上；

### 『阻塞与非阻塞』

阻塞等待的是「**内核数据准备好**」和「**数据从内核态拷⻉到⽤户态**」

这两个过程.先来看看阻塞 I/O，当⽤户程序执⾏ read 线程会被阻塞，⼀直等到内核数据准备好，并把数据从内核缓冲区拷⻉到应⽤程序的缓冲区中，当拷⻉过程完成， read 才会返回。

知道了阻塞 I/O ，来看看⾮阻塞 I/O，⾮阻塞的 read 请求在数据未准备好的情况下⽴即返回，可以继续往下执⾏，此时应⽤程序不断轮询内核，直到数据准备好，内核将数据拷⻉到应⽤程序缓冲区， read 调⽤才可以获取到结果。

注意，这⾥最后⼀次 read 调⽤，获取数据的过程，是⼀个同步的过程，是需要等待的过程。这⾥的同步指的是内核态的数据拷⻉到⽤户程序的缓存区这个过程。举个例⼦，访问管道或 socket 时，如果设置了 O_NONBLOCK 标志，那么就表示使⽤的是⾮阻塞 I/O 的⽅式访问，⽽不做任何设置的话，默认是阻塞 I/O。应⽤程序每次轮询内核的 I/O 是否准备好，感觉有点傻乎乎，因为轮询的过程中，应⽤程序啥也做不了，只是在循环。为了解决这种傻乎乎轮询⽅式，于是 I/O 多路复⽤技术就出来了，如 select、poll，它是通过 I/O 事件分发，当内核数据准备好时，再以事件通知应⽤程序进⾏操作.这个做法⼤⼤改善了应⽤进程对 CPU 的利⽤率，在没有被通知的情况下，应⽤进程可以使⽤ CPU 做其他的事情。下图是使⽤ select I/O 多路复⽤过程。注意， read 获取数据的过程（数据从内核态拷⻉到⽤户态的过程），也是⼀个同步的过程，需要等待：

实际上，⽆论是阻塞 I/O、⾮阻塞 I/O，还是基于⾮阻塞 I/O 的多路复⽤都是同步调⽤。因为它们在 read 调⽤时，内核将数据从内核空间拷⻉到应⽤程序空间，过程都是需要等待的，也
就是说这个过程是同步的，如果内核实现的拷⻉效率不⾼，read 调⽤就会在这个同步过程中等待⽐较⻓的时间.

### 『同步和异步』

⽽真正的异步 I/O 是「内核数据准备好」和「数据从内核态拷⻉到⽤户态」这两个过程都不⽤等待。当我们发起 aio_read 之后，就⽴即返回，内核⾃动将数据从内核空间拷⻉到应⽤程序空间，这个拷⻉过程同样是异步的，内核⾃动完成的，和前⾯的同步操作不⼀样，应⽤程序并不需要主动发起拷⻉动作。

在前⾯我们知道了，I/O 是分为两个过程的：

1. 数据准备的过程
2. 数据从内核空间拷⻉到⽤户进程缓冲区的过程.

阻塞 I/O 会阻塞在「过程 1 」和「过程 2」，⽽⾮阻塞 I/O 和基于⾮阻塞 I/O 的多路复⽤只会阻塞在「过程 2」，所以这三个都可以认为是同步 I/O。
异步 I/O 则不同，「过程 1 」和「过程 2 」都不会阻塞。

阻塞 I/O 好⽐:
去饭堂吃饭，但是饭堂的菜还没做好，然后就⼀直在那⾥等啊等，等了好⻓⼀段时间终于等到饭堂阿姨把菜端了出来（数据准备的过程），但是还得继续等阿姨把菜（内核空间）打到的饭盒⾥（⽤户空间），经历完这两个过程，才可以离开。⾮阻塞 I/O 好⽐，去了饭堂，问阿姨菜做好了没有，阿姨告诉没，就离开了，过⼏⼗分钟，⼜来饭堂问阿姨，阿姨说做好了，于是阿姨帮把菜打到的饭盒⾥，这个过程是得等待的。

基于⾮阻塞的 I/O 多路复⽤好⽐:
去饭堂吃饭，发现有⼀排窗⼝，饭堂阿姨告诉这些窗⼝都还没做好菜，等做好了再通知，于是等啊等（ select 调⽤中），过了⼀会阿姨通知菜做好了，但是不知道哪个窗⼝的菜做好了，⾃⼰看吧。于是只能⼀个⼀个窗⼝去确认，后⾯发现 5 号窗⼝菜做好了，于是让 5 号窗⼝的阿姨帮打菜到饭盒⾥，这个打菜的过程是要等待的，虽然时间不⻓。打完菜后，⾃然就可以离开了。

异步 I/O 好⽐:
让饭堂阿姨将菜做好并把菜打到饭盒⾥后，把饭盒送到⾯前，整个过程都不需要任何等待。
