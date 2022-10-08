---
layout: post
title: 重构代码
subtitle: Composing Methods
tags: [Design Patterns]
---

# 重构代码

## Composing Methods

> 大部分重构都致力于正确组合方法。在大多数情况下，过长的方法是一切邪恶。这些方法中的代码变幻莫测
> 执行逻辑并使该方法极难理解——甚至更难改变。该组中的重构技术简化了方法**删除代码。**

#### `提取方法`

```
func printOwing() {
	printBanner()
	// Print details.
	fmt.Println("name: " + name);	
    fmt.Println("amount: " + getOutstanding());
   
}

func printOwing() {
	printBanner()
	// Print details.
	printDetails(getOutstanding())
  
}
void printDetails(outstanding int, name int) {
	System.out.println("name: " + name);
	System.out.println("amount: " + outstanding)
}
```

解决方案：
将此代码移动到单独的新方法（或函数）和用对方法的调用替换旧代码

为什么要重构？
在方法中找到的行越多，就越难弄清楚该方法的作用。这是造成这种情况的

主要原因：重构
除了消除代码中的粗糙边缘之外，提取方法也是许多其他重构方法中的一个步骤。

好处：
更易读的代码！

如何重构？

创建一个新方法，并以使其纯粹的方式命名。

### `Inline Method`

```
class PizzaDelivery {
	// ...
    int getRating() {
		return moreThanFiveLateDeliveries() ? 2 : 1;
	}
	boolean moreThanFiveLateDeliveries() {
		return numberOfLateDeliveries > 5;
	}
}
```

```
class PizzaDelivery {
	 // ...
	 int getRating() {
		return numberOfLateDeliveries > 5 ? 2 : 1
	}
	boolean moreThanFiveLateDeliveries() {
		 return numberOfLateDeliveries > 5;
	}
}
```

1. Make sure that the method isn’t redefined in subclasses. If the
method is redefined, refrain from this technique.
2. Find all calls to the method. Replace these calls with the con-
tent of the method.
3. Delete the method.

### `Inline Temp`

```

boolean hasDiscount(Order order) {
	double basePrice = order.basePrice();
	return basePrice > 1000;
}

```

```
boolean hasDiscount(Order order) {
	return order.basePrice()> 1000;
}
```

### `Replace Temp with Query`

```
double calculateTotal() {
	double basePrice = quantity * itemPrice;
	 if (basePrice > 1000) {
		return basePrice * 0.95;
	 } else {
		return basePrice * 0.98;
	}
}
```

```
double calculateTotal() {
	double basePrice = quantity * itemPrice;
	 if (basePrice > 1000) {
		return basePrice * 0.95;
	 } else {
		return basePrice * 0.98;
	}
}
double basePrice(){
	return quantity * itemPrice;
}
```

### `Split Temporary Variable`

```
double temp = 2 * (height + width);
System.out.println(temp);
temp = height * width;
System.out.println(temp);
```

```
final double perimeter = 2 * (height + width);
System.out.println(temp);
final double area= height * width;
System.out.println(temp);
```

### `Remove Assignments to Parameters`

```
int discount(int inputVal, int quantity) {
	if (inputVal > 50) {
		inputVal -= 2;
 }
}
```

```
int discount(int inputVal, int quantity) {
	int result = inputVal
	if (inputVal > 50) {
		result -= 2;
 }
}
```

### `Replace Methodnwith Method Object`

```
class Order {
 // ...
 public double price() {
     double primaryBasePrice;
     double secondaryBasePrice;
     double tertiaryBasePrice;
	// Perform long computation.

	}
}
```

```
//订单的价格有多种算法
class Order {
 // ...
 public double price() {
     return new PriceCalculator(this).compute()

	}
}
class PriceCalculator{
	double primaryBasePrice;
     double secondaryBasePrice;
     double tertiaryBasePrice;
	// Perform long computation.
    public PriceCalculator(Order order){
    	//todo
    }
	public double compute() {
		//todo
	}
}
```

### `Substitute Algorithm`

```
String foundPerson(String[] people){
    for (int i = 0; i < people.length; i++) {
		if (people[i].equals("Don")){
			return "Don"
		}
		if (people[i].equals("John")){
			return "John"
		}
        if (people[i].equals("John")){
            return "John"
        }
    }
    return "";	
}
```

```
String foundPerson(String[] people){
	List candidates = Arrays.asList(new String[] {"Don", "John", "Kent"});
    for (int i = 0; i < people.length; i++) {
		if (candidates.contains(people[i])){
			return people[i]
		}
    }
}
```

## `Moving Features between Objects`

> 这些重构技术展示了如何安全地移动函数类，创建新类，并隐藏实现公开访问的详细信息。

- 问题：

  一个方法在另一个类中的使用比在它的类中使用的多自己的班级。

- 解决方案：

  在类中创建一个新方法，使用方法最多，然后将代码从旧方法移到那里。将原始方法的代码转换为对另一个类中新方
  法的引用，否则将其完全删除。

- 问题：一个字段在另一个类中的使用比在它的类中更多自己的班级。

- 解决方案：

  在一个新类中创建一个字段并重定向所有用户

  ```
  type Person struct{
  	name string
  	officeAreaCode string
  	officeNumber string
  }
  ```

  ```
  type Person struct{
  	name string
  	phone TelephoneNumber
  }
  type TelephoneNumber struct{
  	officeAreaCode string
  	officeNumber string
  }
  ```

- 问题：当一个班级做两个班级的工作时，尴尬

- 解决方案：
  相反，创建一个新类并将负责相关功能的字段和方法放入其中

- 问题：一个类几乎什么都不做，也不对任何事情负责，也没有为它计划额外的责任

- 解决方案：

  将所有功能从类移到另一个。

- 问题：客户端从对象 А 的字段或方法中获取对象 B。然后客户端调用对象 B 的一个方法。

- 解决方案：

  在 A 类中创建一个新方法，将调用委托给对象 B。现在客户端不知道或不依赖于 B 类。

- 问题：一个类有太多简单地委托给其他对象的方法。

- 解决办法：

  删除这些方法，强制客户端直接调用方法。创建一个getter 用于从服务器类对象。.用直接调用委托类中的方法替换对服务器类中委托方法的调用。

- 问题：
  实用程序类不包含您需要的方法，并且您不能将该方法添加到类中。

- 解决方案：
  将方法添加到客户端类并传递实用程序类将其作为参数。

  ```
  class Report {
  	// ...
  	void sendReport() {	
  		Date nextDay = new Date(
              previousEnd.getYear(),    
              previousEnd.getMonth(),
              previousEnd.getDate() + 1
  		);    
  		// ...    
  	}
  }
  ```

  ```
  class Report {
  	// ...
  	void sendReport() {	
  		Date newStart = nextDay(previousEnd);   
  		// ...    
  	}
  	private static Date nextDay(Date arg){
  		return new Date(arg.getYear(), arg.getMonth(), arg.getDate() + 1)	
  	}
  }
  ```

#### 组织数据

> 类关联的解开，这使得类更便携和可重用

- 自封装字段

  > 问题：
  > 您使用对类中私有字段的直接访问。
  > 解决方案：
  > 为该字段创建一个 getter 和 setter，并仅使用它们来访问该字段。

- 用对象替换数据值

  > 问题：
  >
  > 一个类（或一组类）包含一个数据字段。该字段有自己的行为和相关数据。
  >
  > 解决方法：
  > 新建一个类，将旧的字段及其行为放在类中，将类的对象存放在原来的类中。

- 将值更改为参考

  >
  > 问题：
  > 您需要用单个对象替换单个类的许多相同实例。
  > 解决方案：
  > 将相同的对象转换为单个参考对象

- 更改对值的引用

  > 问题：
  > 您有一个参考对象太小且很少更改，无法证明管理其生命周期是合理的。
  > 解决方案：
  > 把它变成一个值对象。

- 用对象替换数组

  > 问题：
  > 您有一个包含各种类型数据的数组。
  > 解决方案：
  > 将数组替换为每个元素都有单独的速率字段的对象

- 重复观测数据

  > 问题：
  > 域数据是否存储？
  > 解决方案：
  > 那么最好将数据分离到单独的类中，确保连接和同步

- 将单向关联更改为双向

  > 问题：
  > 你有两个类，每个类都需要使用彼此的关系，但它们之间的关联只是单向的。
  >
  > 解决方案：
  > 将缺少的关联添加到需要它的类中。

- 将双向关联更改为单向

  > 你有一个类之间的双向关联，es，但是其中一个类不使用另一个类的功能。

- 用符号常数替换幻数

  > 问题：
  > 您的代码使用了一个具有特定含义的数字。
  > 解决方案：
  > 将此数字替换为具有人类可读名称的常量，以解释数字的含义。

- 封装字段

  > 问题：
  > 你有一个公共领域。
  > 解决方案：
  > 将字段设为私有并为其创建访问方法

- 封装集合

  > 问题：
  > 一个类包含一个集合字段和一个用于处理集合的简单 getter 和 setter。
  > 解决方案：
  > 将 getter 返回的值设为只读，并创建用于添加/删除集合元素的方法。

- 用类替换类型代码

  > 问题：
  > 一个类有一个包含类型代码的字段。该类型的值不用于操作符条件，也不影响程序的行为。
  > 解决方案：
  > 创建一个新类并使用其对象而不是类型代码值。

- 用子类替换类型代码

  > 问题：
  > 您有一个编码类型，它直接影响每克行为（该字段的值触发条件中的各种代码）。
  > 解决方案：
  > 为编码类型的每个值创建子类。然后将原始类中的相关行为提取到这些子类中。用多态替换控制流代码。

- 用状态/策略替换类型代码

  > 问题：
  > 您有一个影响行为的编码类型，但您不能使用子类来摆脱它。
  > 解决方案：
  > 将类型代码替换为状态对象。如果需要用类型代码替换字段值，则“插入”另一个状态对象

## `Self Encapsulate Field`

自封装

```

class Range {
 	private int low, high;
	boolean includes(int arg) {
		return arg >= low && arg <= high;
	}
}
```

```

class Range {
 	private int low, high;
	boolean includes(int arg) {
		return arg >= getLow()&& arg <= getHigh();
	}
	int getLow(){
		return low;
	}
	int getHigh(){
		return high;
	}
	
}
```

用对象替换数据值

```
type Order struct{
	Customer string
}
```

```
type Order struct{
	Customer Customer
}
type Customer struct{
	Name string
}
```

通过用对象替换数据值，我们有了一个原始字段（数字、字符串等），由于程序的增长，它不再那么简单，现在有了相关的数据和行为。一方面,这些领域本身并没有什么可怕的。但是，这个字段和行为系列可以同时存在于多个类中，从而创建重复的代码。

将值更改为引用

```
type Order struct{
	Customer Customer
}
type Customer struct{
	Name string
}
```

```
type Order struct{
	Customer *Customer
}
type Customer struct{
	Name string
}
```

将引用更改为值

```
type Customer struct{
	Currency *Currency 
}
type Currency struct{
	Code  string
}
```

```
type Customer struct{
	Currency Currency 
}
type Currency struct{
	Code  string
}
```

将数组替换为对象

```
  String[] row = new String[2];
  row[0] = "Liverpool";
  row[1] = "15";
```

```
  Performance row = new Performance();
  row.setName("Liverpool");
  row.setWins("15");
 nn
```

重复观察数据

```
type IntervalWindow struct{
	Textstart sting
	Textend string
	legth int
}
func (I *IntervalWindow) CalculateLegth(){
	
}
func (I *IntervalWindow) CalculateEnd(){
	
}
func (I *IntervalWindow) CalculateStart(){
	
}


```

```

type IntervalWindow struct{
	 Interval  Interval 
} 
type Interval struct{
	Start sting
	End string
	legth int
}

func (I *Interval) CalculateLegth(){
	
}
func (I *Interval) CalculateEnd(){
	
}
func (I *Interval) CalculateStart(){
	
}
```

单向关系变为双向

```
type Order struct{
	Customer *Customer
}
type Customer struct{
	Name string
}
```

```
type Order struct{
	Customer *Customer
}
type Customer struct{
	Order *Order
	Name string
}
```

变双向为单向

```
type Order struct{
	Customer *Customer
}
type Customer struct{
	Order *Order
	Name string
}
```

```
type Order struct{
	Customer *Customer
}
type Customer struct{
	Name string
}
```

替换魔法 编号和符号常量

```

double potentialEnergy(double mass, double height) {
	return mass * height * 9.81;
}
```

```

static final double GRAVITATIONAL_CONSTANT = 9.81;

double potentialEnergy(double mass, double height) {
	return mass * height * GRAVITATIONAL_CONSTANT;
}
```

封装数组

```
type Teacher struct{
	PersonList  []Person
}
type Person struct{
	Name string
}

```

```
type Teacher struct{
	PersonList  []Person
}
func (t *Teacher) Getter(index int){
	return t.PersonList[index]
}
func (t *Teacher) Setter(p Person , i int){
	t.PersonList[i]=p
}
type Person struct{
	Name string
}
```

类型

```
type Teacher struct{
	O  int
	B  int
	C  int
	AB  int
}

```

```
type Teacher struct{
	O Bloodgroup
	B Bloodgroup
	C Bloodgroup
	AB Bloodgroup
}
type Bloodgroup int
```

有些字段无法被验证，由 IDE 检查类型。

用子类替换类型代码

```
type Empployee struct{
	engineer int
	salesman int
}
```

```
type Empployee struct{
	Engineer Engineer 
	Salesman Salesman 
}
type Engineer struct{
	Id int
}
type Salesman  struct{
	Id int
}
```

用状态代替类型代码

```
type Empployee struct{
	Engineer Engineer 
	Salesman Salesman 
}
type Engineer struct{
	Id int
}
type Salesman  struct{
	Id int
}
```

```
type Empployee  EmpployeeType
type EmpployeeType struct{
	Engineer Engineer 
	Salesman Salesman 
}
type Engineer struct{
	Id int
}
type Salesman  struct{
	Id int
}
```

用字段替换子类

```
type Person struct{
	 Male  *Male
	 FeMale *FeMale
}
func (p *Person)GetCode(){
	
}

type Male struct{
}
func (p *Male)GetName()string{
	return "M"
}

type  FeMale struct{
}
func (p *FeMale)GetName()string{
	return "F"
}
```

```
type Person struct{
	 code string
	 
}
func (p *Person)GetCode(){
	
}

```

## `Simplifying Conditional Expressions`

> 问题：
> 您有一个复杂的条件（ if‑then / else或switch ）。
> 解决方案：
> 将条件的复杂部分分解为单独的方法：条件，然后和其他。

> 问题：
> 您有多个导致相同结果或操作的条件。
> 解决方案：
> 将所有这些条件合并到一个表达式中

> 问题：
> 在所有分支中都可以找到相同的代码有条件的。
> 解决方案：
> 将代码移到条件之外。

> 问题：
> 您有一个用作控件的布尔变量多个布尔表达式的标志。
> 解决方案：
> 代替变量，使用break, continue 和return。

> 问题：
> 你有一组嵌套条件，这很难,,来确定代码执行的正常流程。
> 解决方案：
> 将所有特殊检查和边缘情况隔离到单独的子句中，并将它们放在主要检查之前。理想的盟友，你应该有一个“平坦”的条件列表，一个接一个,另一个。

> 问题：
> 您有一个执行各种操作的条件,取决于对象类型或属性。
> 解决方案：
> 创建与**条件分支匹配的子类**。在其中，创建一个共享方法并从条件的相应分支。然后更换
> 带有相关方法调用的条件。结果是正确的实现将通过多态性获得，具体取决于对象类

> 问题：
> 由于某些方法返回null而不是 real,对象，您在代码中对null进行了许多检查。
> 解决方案：
> 而不是null, 返回一个展示的空对象.默认行为。

>
> 问题：
> 要让一部分代码正常工作，某些条件或值必须为真。
> 解决方案：
> 用特定的断言替换这些假设检查。

## `Decompose Conditional`

```
if (date.before(SUMMER_START) || date.after(SUMMER_END)) {
	charge = quantity * winterRate + winterServiceCharge;
}
else {
		charge = quantity * summerRate;
	}
```

```
if (isSummer(date)) {
	charge = summerCharge(quantity);
}
 else {
	charge = winterCharge(quantity);
}
```



```
double disabilityAmount() {
if (seniority < 2) {
	return 0;
 }
 if (monthsDisabled > 12) {
	return 0;
 }
 if (isPartTime) {
	return 0;
 }
	// Compute the disability amount.
	 // ...nn
}
```

```
double disabilityAmount() {
  if (isNotEligableForDisability()) {
	return 0;	
  nn}
	// Compute the disability amount.
	// ...
}
```



```
if (isSpecialDeal()) {
	 total = price * 0.95;
	 send();
}
else {
	total = price * 0.98;
	send();
}
```

```

if (isSpecialDeal()) {
	total = price * 0.95;
}
	else {
	total = price * 0.98;
 }	
 nsend();
```



```
break: stops loop
continue: stops execution of the current loop branch and
goes to check the loop conditions in the next iteration
return: stops execution of the entire function and returns its
result if given in the operator
```



## `Replace Nested Conditional with Guard Clauses`

```
	
public double getPayAmount() {
	 double result;
	if (isDead){
		result = deadAmount();
	}else {
		if (isSeparated){
			result = separatedAmount();
		}else {
			if (isRetired){
	r			esult = retiredAmount();
			}else{
				result = normalPayAmount();
			}
		}
	}
	return result
}
```

```
public double getPayAmount() {
	if (isDead){
		return deadAmount();
	 }
	if (isSeparated){
		return separatedAmount();
	}
	if (isRetired){
		return retiredAmount();
	}
	 return normalPayAmount();
}
```



## `Replace Conditional with Polymorphism`

```
class Bird {
 // ...
 double getSpeed() {
 	switch (type) {
		case EUROPEAN:
			return getBaseSpeed();
		case AFRICAN:
			return getBaseSpeed() - getLoadFactor() * numberOfCoconuts;
		case NORWEGIAN_BLUE:
			return (isNailed) ? 0 : getBaseSpeed(voltage);

	 }
	throw new RuntimeException("Should be unreachable");
  }
}
```

