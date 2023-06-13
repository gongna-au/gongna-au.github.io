---
layout: post
title: Java/Go/C++
subtitle:
tags: [Java]
comments: true
---


### Slice

```java
public class Main {
    public static void main(String[] args) {
        List<Integer> slice = new ArrayList<>();
        slice.add(1);
        slice.add(2);
        slice.add(3);
        slice.add(4);
        slice.add(5);
        slice = slice.subList(2, 4);
        slice.addAll(List.of(1, 2, 3));
        System.out.println(slice);
    }
}

```


```go
package main

import "fmt"

func main() {
    slice := []int{}
    slice = append(slice, 1)
    slice = append(slice, 2)
    slice = append(slice, 3)
    slice = append(slice, 4)
    slice = append(slice, 5)
    slice = slice[2:4]
    slice = append(slice[:len(slice)-1], []int{1, 2, 3}...)
    fmt.Println(slice)
}

```

```c++
#include <iostream>  
#include <vector>
#include <iterator>

int main() { 
     std::vector<int> slice;
     slice.push_back(1);
     slice.push_back(2);
     slice.push_back(3);
     slice.push_back(4);
     slice.push_back(5);

     std::vector<int>::iterator it1 = slice.begin();   
     std::advance(it1, 2);                
     slice.erase(it1, slice.end());                 

     std::vector<int> append_values;
     append_values.push_back(1);
     append_values.push_back(2); 
     append_values.push_back(3);

     slice.insert(slice.end(), append_values.begin(), append_values.end());               
                                                                                    
     for (std::vector<int>::iterator it = slice.begin(); it != slice.end(); it++) {
         std::cout << *it << ' ';
     }
     std::cout << std::endl;
     return 0;
}
```

### Map

```java
import java.util.Map;
import java.util.HashMap; 

public class Main {
    public static void main(String[] args) {
        Map<Integer, String> hashMap = new HashMap<>();
        hashMap.put(1, "one");   
        hashMap.put(2, "two");
        hashMap.put(3, "three");
        
        System.out.println(hashMap.get(1));  // 输出 one
        
        hashMap.remove(3);
        
        for (Map.Entry<Integer, String> entry : hashMap.entrySet()) {
            int key = entry.getKey();  
            String value = entry.getValue();
            System.out.println(key + " : " + value);
        }
    }
}
```

```go
package main

import "fmt"

func main() {
    hashMap := make(map[int]string)
    
    hashMap[1] = "one"
    hashMap[2] = "two"    
    hashMap[3] = "three"
    
    fmt.Println(hashMap[1])  // 输出 one
    
    delete(hashMap, 3)
    
    for key, value := range hashMap {
        fmt.Printf("%d : %s\n", key, value)
    }  
}
```


```c++
#include <iostream>
#include <unordered_map>
#include <string>

using namespace std;

using HashMap = unordered_map<int, string>;

int main() {
    HashMap hashMap;
    hashMap[1] = "one";
    hashMap[2] = "two";
    hashMap[3] = "three";
    cout << hashMap[1] << endl;
    hashMap.erase(3);
    return 0;
}
```



### 接口

```java
List<Integer> list = new ArrayList<Integer>();
```

> List 接口是 Java 集合框架中的一部分，它定义了一个有序的集合，其中的元素可以重复。List 接口中定义了一些常用的方法，例如 add、get、remove 等等。而 ArrayList 类则是 List 接口的一个具体实现，它使用数组来实现列表，可以动态扩展和收缩，可以在列表的任意位置随机访问元素，是 Java 中最常用的列表实现类之一。



```java
Map<Integer,String> map = new HashMap<Integer>();
```

> 使用了 HashMap 来创建一个 Map<Integer, String> 对象。这里的 HashMap 是一个具体的实现类，实现了 Map 接口，可以用来创建 Map 对象。需要注意的是，HashMap 是一种基于哈希表的实现类，可以快速地查找键对应的值，是最常用的 Map 实现类之一。



### Java的Byte[] 和Go的[]byte

```java
byte[] bytes = {72, 101, 108, 108, 111, 33};
String str = new String(bytes);
System.out.println(str); // 输出 Hello!
```
或者
```java
public class Main {
    public static void main(String[] args) {
        byte[] bytes = {'H', 'e', 'l', 'l', 'o', '!'};
        String str = new String(bytes);
        System.out.println(str); // 输出 Hello!
    }
}
```

```go
s:=[]byte{'h','e','l','l','o'}
fmt.Println(String(s))
```


### Static/Const 

static 可以用来声明一个静态变量,也可以用来声明一个静态函数。

静态变量（`static`）和常量（`const`）在 C++ 中具有不同的含义和使用场景。

1-`static` 关键字：
- **在全局变量和函数中**：`static` 使全局变量和函数的范围局限于定义它们的文件。换句话说，一个在文件 A 中定义的`static` 变量或函数不能在文件 B 中被访问。
- **在局部变量中**：`static` 使局部变量的生命周期在程序运行期间持续存在，而不是在其所在的函数或代码块结束时消亡。也就是说，它们的值会在多次函数调用中保持不变。
- **在类中**：`static` 使类的成员不再依赖于特定的类的实例。这意味着 `static` 成员只存在一份，被所有的类实例共享。

2. `const` 关键字：
- **变量**：`const` 使得变量的值在声明后不可修改。尝试修改 `const` 变量的值会导致编译错误。
- **函数**：在函数声明中，`const` 关键字表示该函数不会修改它的对象的状态。这对于理解对象的状态在何时被修改很重要。

没有 `static` 或 `const` 声明的变量和函数有默认的可见性和生命周期（在定义它们的范围内），并且它们的值可以随时被修改。


```c++
class A {
    int num;
};
  
void f(A obj) { // 普通函数,可以修改 obj 
    obj.num = 10; 
}

void f(const A obj) { // 常量函数,不会修改 obj
     obj.num = 10; // 编译错误! obj 是 const 的
}
```

```c++
const int num = 10; // num 是一个常量,不能修改
num = 20; // 编译错误!
```

```text
char *const p = "ABCD";
这里,*const 表明p是一个const(常量)指针。
*p = 'X'; // 修改第一个字符为'X'
p[1] = 'Y'; // 修改第二个字符为'Y'
可以修改p指向的内容
但是不能：
p = "1234"; // 错误!不能修改const指针本身
```

> 简言之，可以修改现在有两个内存块A B ，p指向A，如果是 char * const p ，那么意味你可以把A内存块的内容改为C，但是不能让p指向新的内存块B


> 相反的const int*  p 可以让p指向新的内存块B，但是不能更改p指向的内存块的内容。

> 在函数调用过程：值会从上一次结束时的值开始
```c++
#include <stdio.h>

void func() {
    static int count = 0;
    count++;
    printf("Count: %d\n", count);
}

int main() {
    func();
    func();
    func();
    return 0;
}
输出：
Count: 1
Count: 2
Count: 3
```


> 在模块内（但在函数体外），**静态的变量**可以被模块内所用函数访问，但不能被模块外其它函数访问。它是一个本地的全局变量。

```c++
#include <stdio.h>

void func1();
void func2();

static int count = 0;

int main() {
    func1();
    func2();
    func1();
    return 0;
}

void func1() {
    count++;
    printf("Count in func1: %d\n", count);
}

void func2() {
    count += 2;
    printf("Count in func2: %d\n", count);
}
```


> 在模块内，一个被声明为静态的函数只可被这一模块内的其它函数调用。那就是，这个函数被限制在声明它的模块的本地范围内使用。
```c++
#include <stdio.h>

static void func1();
static void func2();

int main() {
    func1();
    func2();
    return 0;
}

static void func1() {
    printf("This is func1\n");
}

static void func2() {
    printf("This is func2\n");
}
```


