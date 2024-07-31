---
layout: post
title: 软件测试-IBM Rational Functional Tester
subtitle:
tags: [软件测试]
comments: true
---

> https://www.ibm.com/docs/en/SSJMXE_10.5.1/docs/pdfs/Rational_Functional_Tester.pdf

### Lab O Become Familiar with the Sample Application 熟悉示例应用程序

### Lab 1 Getting Started with IBM Rational Functional Tester 开始使用 IBM Rational Functional Tester

### Lab 2 Recording a Script 录制脚本


1. 打开 IBM Rational Functional Tester 软件并创建一个新项目。
2. 在项目中创建一个新的测试用例。
3. 选择要测试的应用程序并启动它。
4. 在 Rational Functional Tester 中单击“记录”按钮来开始录制脚本。
5. 在应用程序中模拟用户的行为，例如单击按钮，填写表单等。
6. 在完成操作后停止录制。
7. 对录制的脚本进行编辑和修改，使其符合需要。
8. 运行测试脚本并查看结果。

请注意，在录制期间，确保 Rational Functional Tester 已正确配置以与应用程序集成，并且您已打开了“扫描到应用程序”选项。



### Lab 3 Playing Back a Script and Viewing Results 回放脚本并查看结果


1. 打开前面创建的项目和测试用例。
   您需要打开 IBM Rational Functional Tester，并选择相应的项目和测试用例。

2. 点击“回放”按钮，并选择“Java Test Script playback”选项。
   在Rational Functional Tester界面的工具栏上，您可以看到“回放”按钮。单击该按钮并选择“Java Test Script playback”。

3. 运行脚本时，请确保Rational Functional Tester已正确配置以与应用程序进行集成，并且已打开“扫描到应用程序”选项。
   在运行测试脚本之前，需要确认 Rational Functional Tester 已经完成与被测应用程序的集成，并且已在“全局设置” > “Java” 中打开扫描到应用程序选项。

4. 当脚本运行完成后，您将看到执行结果（通过绿色/红色标记显示）。您还可以查看详细的执行结果和报告。
   在测试脚本运行完成后，您会在 Rational Functional Tester 界面中看到绿色或红色的标记，表示测试脚本是否成功执行完成。如果您希望查看更加详细的执行结果和报告，可以在 Rational Functional Tester 的测试结果视图中查看。

### Lab 4 Extending Scripts 扩展脚本


1. 在 IBM Rational Functional Tester 中，创建一个新项目。
   打开Rational Functional Tester并创建一个新项目。在“文件”菜单中选择“新建”、“项目”，然后根据需要指定项目名称和位置。

2. 在该项目下，创建一个新的Java测试脚本。
   接下来，在新项目下创建一个新的 Java 测试脚本。在“测试”菜单中选择“新建”、“Java 测试脚本”。给这个测试脚本命名，并确认它被正确保存在您的项目中。

3. 向测试脚本中添加一些基本的测试步骤，如启动应用程序、输入文本并单击按钮等。
   首先，您需要向测试脚本中添加一些基本的测试步骤，例如启动应用程序，输入文本并单击按钮等。这些测试步骤将成为测试脚本的基础，以便在此之后进行扩展。

4. 使用 Rational Functional Tester 的记录功能来记录一些新的测试步骤。
   接下来，使用 Rational Functional Tester 的记录功能来记录一些新的测试步骤。您可以通过单击“记录”按钮来启动记录模式，然后执行您想要记录的测试步骤。一旦完成，单击“停止”按钮以停止记录模式。

5. 编辑测试脚本以包含新记录的测试步骤，并根据需要进行调整和修改。
   现在，您可以编辑测试脚本以包含新记录的测试步骤。可能需要对新测试步骤进行一些修改和调整，以确保它们与先前的测试步骤协同工作。

6. 使用 Java 编写扩展代码，例如将测试数据读入文件、使用迭代循环运行测试等。
   接下来，您可以使用 Java 编写扩展代码，例如将测试数据读入文件、使用迭代循环运行测试等。这些扩展代码将使测试脚本更加灵活和可复用。

7. 运行测试，查看结果并进行优化。
   最后，运行您的测试，并查看其结果。如果有任何失败或错误，请查看测试日志和报告，并根据需要进行调整和优化。



### Lab 5 Using Test Object Maps 使用测试对象映射


1. 在 IBM Rational Functional Tester 中，创建一个新项目。
   打开Rational Functional Tester并创建一个新项目。在“文件”菜单中选择“新建”、“项目”，然后根据需要指定项目名称和位置。

2. 在该项目下，创建一个新的Java测试脚本。
   接下来，在新项目下创建一个新的 Java 测试脚本。在“测试”菜单中选择“新建”、“Java 测试脚本”。给这个测试脚本命名，并确认它被正确保存在您的项目中。

3. 启动被测试应用程序，并让 Rational Functional Tester 检测该应用程序。
   首先，启动要测试的应用程序，并确保 Rational Functional Tester 已经检测到该应用程序。您可以通过单击“Window”菜单中的“Functional Test Project Map”选项卡来查看是否正确检测到了应用程序。

4. 创建一个新的测试对象映射。
   接下来，创建一个新的测试对象映射。在“测试”菜单中选择“新建”、“对象映射”，并将其命名为适当的名称。

5. 输入测试对象信息。
   在测试对象映射编辑器中，输入有关测试对象的信息。该信息包括对象的类型、属性和值等。

6. 使用测试对象映射中的对象来创建测试步骤。
   接下来，使用测试对象映射中的对象来创建您的测试步骤。您可以单击测试对象映射编辑器中的“插入”按钮，并选择要在测试脚本中创建的测试步骤类型。

7. 运行测试并查看结果。
   最后，运行您的测试，并查看其结果。如果有任何失败或错误，请查看测试日志和报告，并根据需要进行调整和优化。



### Lab 6 Managing Object Recognition 管理对象识别


1. 启动IBM Rational Functional Tester并打开被测试的Java应用程序。

2. 在测试脚本窗口中创建一个新Java测试脚本文件。

3. 让Rational Functional Tester自动识别应用程序中的对象。

4. 确认对象的正确识别。

5. 调整对象属性以进行更准确的对象识别。您可以添加、删除或编辑对象属性。

6. 执行测试，并查看测试结果。如果有任何失败或错误，请查看测试日志和报告，并根据需要进行调整和优化。

7. 测试完成后，保存并关闭测试脚本文件。



### Lab 7 Creating Data-Driven Tests 创建数据驱动测试




1. 启动IBM Rational Functional Tester并打开被测试的应用程序。

   - 在启动 Rational Functional Tester之前，请确保已经安装了Java JDK和Eclipse IDE。
   
   - 打开Rational Functional Tester，单击File > New > Project，选择Java工程模板，在Project Name文本框中输入项目名称并单击Finish。

   - 单击File > New > Test Script，选择Java模板创建一个新的Java测试脚本文件。在弹出对话框中，选择Test Object Map（TOM）作为对象库。然后单击Next。

   - 选择要测试的应用程序，通过点击“Start Application”按钮来启动该应用程序。

2. 准备测试数据。

   - 您可以使用Excel电子表格或文本文件来存储测试数据。将测试数据保存在文件中，并确定它们可供读取。

3. 将测试数据导入Rational Functional Tester。

   - 在Rational Functional Tester中，单击File > Import > Data-Driven Test，选择要使用的测试数据文件类型（Excel、文本等）。设置好文件路径等参数，然后单击Finish来导入测试数据。

4. 在测试脚本中使用测试数据。

   - 将测试数据的值分配给变量，使用for循环、if语句和其他控制流语句来操作这些变量并执行测试。

   - 在测试脚本窗口中，找到导入的测试数据文件，在“Data Driven Scripts”视图下双击该文件，然后选择要使用的数据集并单击OK。

5. 执行测试，并查看测试结果。

   - 单击Run按钮或在菜单栏中单击Run > Run As > Rational Functional Tester Test。测试运行时，测试脚本将使用指定的测试数据来执行测试步骤。

6. 测试完成后，保存并关闭测试脚本文件。

   - 在测试脚本窗口中，单击File > Save to SCM。在弹出的对话框中，输入注释并单击OK。



### IBM RFT 8.5 

#### Part 1 - Introduction

在IBM RFT 8.5 - Part 1介绍中，您将了解以下内容：

1. 什么是RFT？
2. RFT的特点和功能。
3. 如何创建一个新项目。

##### 什么是RFT？

RFT是IBM Rational系列测试工具之一，它提供了一个完整的自动化测试解决方案，可以用于各种类型的应用程序测试，包括Web应用程序，Java应用程序，.NET应用程序等。 

RFT使用脚本语言，称为Rational Functional Tester Script（RFT脚本），该脚本语言基于Java编写，带有丰富的API和内置操作，可用于执行各种测试操作，例如模拟用户交互和检查应用程序状态。

##### RFT的特点和功能

RFT具有以下功能：

- 自动化GUI测试
- 支持多种技术，例如Java，HTML，Net，Oracle等
- 可以同时运行多个自动化测试
- 提供记录和回放测试脚本的功能
- 提供强大的测试结果分析功能
- 直接集成到Eclipse和IBM Rational Application Developer (RAD) 中。

RFT脚本支持不同的任务，包括自动化测试执行、数据操纵、异步任务调度、异常处理等。 此外，RFT还提供了可扩展性和定制性，可以根据特定需求创建自定义插件和微件。

##### 如何创建新项目

要在RFT中创建新项目，请按照以下步骤操作：

1. 打开IBM Rational Functional Tester。
2. 单击“文件”>“新建”>“Rational Functional Tester项目”。
3. 在“新项目”对话框中，为项目命名，选择测试类型并指定存储位置。
4. 按照向导的指示完成项目设置。

完成上述步骤后，您可以开始使用RFT创建和管理自己的自动化测试。 

总之，IBM RFT是一款功能强大的自动化测试工具，提供了丰富的API和内置功能，可用于各种类型的应用程序测试。 它易于使用，可扩展和定制，并且可以良好地与Eclipse和IBM Rational Application Developer集成。

#### Part 2 - Objects, Domains

在IBM RFT 8.5 - Part 2中，我们将了解以下内容：

1. 对象和对象识别器
2. 域和域测试对象
3. 如何创建自定义域和测试对象

##### 对象和对象识别器

在RFT中，一个对象是指应用程序中的一个控件或部分。可以通过对象识别器在RFT中识别对象。对象识别器使用属性来标识对象，并将其表示为测试脚本中的代码。

在RFT中，有几种不同类型的对象识别器，例如基于图像、基于文本、基于注释等。它们可以根据需要进行切换，以识别各种对象和控件。

##### 域和域测试对象

在RFT中，域是包含相似对象的组。例如，对于Web应用程序，可能会有一个“表单”域，其中包含所有表单元素，如文本框、复选框、按钮等。

每个域都具有一个特定的测试对象。这些对象提供了与域相关的属性和操作。例如，在表单域中，您可以使用测试对象获取表格中的数据或输入表单字段。

##### 如何创建自定义域和测试对象

要创建自定义域和测试对象，请按照以下步骤操作：

1. 在RFT中打开项目。
2. 单击“域”视图中的“新建”按钮。
3. 选择要创建的域类型，然后按照向导的指示完成域设置。
4. 在“测试对象”视图中，单击“新建”按钮。
5. 选择要创建的测试对象类型，输入名称并指定父域。按照向导的指示完成测试对象设置。

完成上述步骤后，您可以使用自定义域和测试对象来执行特定类型的自动化测试。

总之，在RFT中，对象是应用程序中的控件或部分，可以通过对象识别器进行识别。 域是包含相似对象的组，并且每个域都有一个特定的测试对象，提供与域相关的属性和操作。 您可以根据需要创建自定义域和测试对象。



#### Part 3 - Record script

好！我很乐意帮助您使用 IBM RFT 8.5 录制脚本。以下是步骤：

1. 打开 IBM RFT 并通过选择 File > New > Test Script 创建一个新的测试脚本。
2. 在对象查找器窗格中，选择要记录交互的应用程序或网页。
3. 单击工具栏中的 Record 按钮或转到 Actions > Record Testcase。
4. 执行您在应用程序或网页上的所需操作，例如单击按钮或输入文本到字段中。当您这样做时，RFT 将根据您的交互自动生成代码到测试脚本编辑器中。
5. 完成您的交互后，单击工具栏中的停止录制按钮或转到 Actions > Stop Recording。



#### Part 4 - Private Object Map

##### 什么是私有对象映射？
私有对象映射是IBM RFT工具的一个重要特性。如果您使用RFT开发自动化测试脚本，则必须掌握私有对象映射。

私有对象映射是一种机制，可以为测试对象定义逻辑名称，从而更方便地识别和操作这些对象。它与基于层次结构的对象查找不同，后者使用对象的物理位置来查找和标识对象。

在私有对象映射中，您可以为每个测试对象定义一个唯一的逻辑名称，然后在测试脚本中使用该名称来识别和访问该对象。这样做可以使测试脚本更加可读、易于维护和可靠。

##### 如何创建私有对象映射？
要创建私有对象映射，请按照以下步骤操作：

1. 打开应用程序或Web页面，并启动RFT。
2. 在RFT中创建一个新的测试脚本。
3. 选择菜单栏中的“对象”>“添加到私有对象映射”选项。该选项将打开私有对象映射编辑器。
4. 在私有对象映射编辑器中，单击“添加”按钮，并指定测试对象的逻辑名称和物理位置。可以使用“捕获”按钮来记录测试对象的物理位置。
5. 添加所有需要使用的测试对象。
6. 单击“保存并关闭”按钮，以保存私有对象映射更改并关闭编辑器。

##### 如何在测试脚本中使用私有对象映射？
完成私有对象映射的创建后，您可以在测试脚本中使用定义的逻辑名称来引用测试对象。例如，您可以使用以下代码来单击名为“loginButton”的按钮：

```
TestObject loginButton = find("loginButton");
((GuiTestObject)loginButton).click();
```

在这里，“loginButton”是您为该按钮定义的逻辑名称。 `find`方法使用该名称查找测试对象，并将其返回给`testObject`变量中。

注意：在使用RFT时，“find”方法是一种非常常见的技术。它允许您根据某些标准（如逻辑名称或属性）查找测试对象，并在测试脚本中使用这些对象执行操作。

##### 总结
私有对象映射是IBM RFT工具的一个重要特性，可帮助测试人员更方便地识别和操作测试对象。通过在测试脚本中使用逻辑名称而不是物理位置来引用测试对象，可以使测试脚本更加可读、易于维护和可靠。

#### Part 5 - Update recognition properties


##### 什么是识别属性？
在创建RFT测试脚本时，您需要指定如何识别应用程序或Web页面中的各个测试对象。例如，如果您要单击登录按钮，则必须先找到该按钮并使用适当的方法调用来单击它。

要找到测试对象，RFT使用一组称为“识别属性”的属性。这些属性描述了测试对象的外观、位置和其他特征。

##### 何时需要更新识别属性？
有时，应用程序或Web页面的更改可能会影响测试对象的外观或位置，从而导致RFT无法正确识别该对象。在这种情况下，您需要更新识别属性以反映更改，以确保RFT可以准确地识别该对象。

此外，针对某些对象，您可能希望添加自定义识别属性，以便更容易地识别它们。

##### 如何更新识别属性？
要更新测试对象的识别属性，请按照以下步骤操作：

1. 在RFT中打开测试脚本并启动应用程序或Web页面。
2. 使用Object Finder工具（Object Finder Tool）选择要更改的测试对象。
3. 完全标识测试对象后，右键单击它，并选择“更改对象属性”选项。
4. 在“对象属性”窗口中，查看测试对象的当前识别属性设置。
5. 如果需要更新现有识别属性，请单击“编辑”按钮并更改属性的值。
6. 如果需要添加自定义识别属性，请单击“添加”按钮并指定自定义属性的名称和值。
7. 单击“确定”按钮以保存更改，并关闭“对象属性”窗口。

完成这些步骤后，RFT将使用新的识别属性来查找和识别测试对象。

##### 总结
在创建RFT测试脚本时，您需要指定如何识别应用程序或Web页面中的各个测试对象。要找到测试对象，RFT使用一组称为“识别属性”的属性。如果应用程序或Web页面更改导致RFT无法正确识别测试对象，则可以更新或添加识别属性以解决此问题。

#### Part 6 - Simple script

##### 创建一个新测试脚本
1. 打开IBM RFT。
2. 在“RFT工作台”窗口中，选择“文件”>“新建”>“测试脚本”。
3. 在“新建测试脚本”向导中，指定测试脚本名称和项目位置，然后单击“下一步”。
4. 在“应用程序描述”页面上，选择要测试的应用程序类型，并填写应用程序的文件路径或URL。单击“下一步”。
5. 在“测试环境配置”页面上，接受默认值并单击“下一步”。
6. 在“测试域配置”页面上，接受默认值并单击“下一步”。
7. 在“测试对象映射”页面上，接受默认值并单击“下一步”。
8. 在“测试脚本生成器”页面上，选择“手动录制所有操作”选项并单击“完成”。

##### 录制测试脚本步骤
在创建了新的测试脚本之后，您可以执行以下步骤:

1. 在RFT的“测试脚本”窗口中，单击“录制”按钮以开始录制脚本。
2. 启动要测试的应用程序，并通过应用程序执行所需的操作。
3. 在完成测试后，在RFT的“测试脚本”窗口中单击“停止”按钮以停止录制。
4. 在“测试脚本”窗口中，可以查看和编辑刚刚录制的脚本内容。

##### 运行测试脚本
1. 确保要测试的应用程序已在本地或远程计算机上启动。
2. 在RFT的“测试脚本”窗口中选择要运行的脚本。
3. 单击“运行”按钮以运行所选的脚本。
4. RFT将自动执行脚本中记录的操作并生成报告。

##### 总结
创建简单的RFT测试脚本需要执行以下步骤：创建一个新的测试脚本、录制测试脚本步骤和运行测试脚本。您可以使用RFT的图形界面来完成这些任务，并根据需要对录制的脚本进行编辑和修改。

#### Part 7 - Datapool basics

##### 什么是数据池？
数据池是指一组测试数据的集合。在测试期间使用数据池可以帮助您模拟各种不同的测试环境和情况。

##### 如何创建数据池？
要创建一个数据池，请按照以下步骤操作：

1. 打开RFT并创建一个新的测试脚本。
2. 在“测试脚本”窗口中，单击“数据池”按钮。
3. 在“数据池”窗口中，单击“创建新数据池”按钮。
4. 指定数据池名称、数据类型和字段数，并单击“确定”按钮。

##### 如何添加测试数据到数据池中？
要将测试数据添加到数据池中，请执行以下操作：

1. 在“数据池”窗口中选择要添加数据的数据池。
2. 单击“添加记录”按钮。
3. 在“添加数据记录”对话框中，输入所需的测试数据。
4. 单击“确定”按钮以添加数据。

您可以重复这些步骤来添加多个测试数据到数据池中。

##### 如何将数据池与测试脚本关联？
要将数据池与测试脚本关联，请执行以下操作：

1. 打开测试脚本并转到需要使用数据池的测试对象。
2. 在测试对象上右键单击并选择“属性...”。
3. 在“属性”窗口中，选择“数据驱动”选项卡。
4. 选择一个数据池并指定要使用的数据列。
5. 单击“确定”按钮以保存更改。

现在您的测试脚本已配置为使用所选数据池中的测试数据来执行测试操作。

##### 总结
在RFT中，使用数据池可帮助您模拟各种不同的测试情况和环境。要创建数据池，请打开RFT并创建一个新的测试脚本。然后，在“数据池”窗口中添加测试数据并将其与测试脚本关联。这样，您的测试脚本就可以使用数据池中的测试数据来执行测试操作。

#### Part 8 - Datapool management


##### 如何修改数据池？
要修改数据池，请执行以下操作：

1. 打开RFT并转到“测试脚本”窗口。
2. 在左侧窗格中单击“数据池”按钮。
3. 在“数据池”窗口中，选择要修改的数据池。
4. 单击“编辑数据池”按钮。
5. 对数据池进行所需的更改，例如添加或删除数据行。
6. 单击“确定”按钮以保存更改。

##### 如何导入和导出数据池？
要导入数据池，请执行以下操作：

1. 在“数据池”窗口中，单击“导入”按钮。
2. 浏览到包含要导入的数据池文件的位置，并输入文件名。
3. 单击“打开”按钮以开始导入过程。
4. 根据需要调整数据池设置，例如字段分隔符和日期格式。
5. 确认导入选项，然后单击“完成”按钮。

要导出数据池，请执行以下操作：

1. 在“数据池”窗口中，选择要导出的数据池。
2. 单击“导出”按钮。
3. 浏览到要将数据池文件保存的位置，并输入文件名。
4. 根据需要调整导出选项，例如字段分隔符和日期格式。
5. 单击“确定”按钮以开始导出过程。

##### 如何删除数据池？
要删除数据池，请执行以下操作：

1. 在“数据池”窗口中，选择要删除的数据池。
2. 单击“删除数据池”按钮。
3. 在提示框中单击“是”以确认删除操作。

注意：删除数据池后，其中的所有测试数据都将永久丢失。请谨慎操作。

##### 总结
在RFT中，您可以随时修改、导入、导出和删除数据池。要修改数据池，打开“数据池”窗口并对数据进行更改。要导入或导出数据池，请单击相应的按钮并进一步指定有关选项。要删除数据池，请选择其名称并单击“删除数据池”按钮。


#### Part 9 - Run script with multiple data

要在IBM RFT 8.5中使用多个数据集运行脚本，可以使用数据驱动测试功能。以下是要遵循的步骤：

1. 创建新的测试脚本或打开现有脚本。

2. 单击工具栏上的“数据驱动测试”按钮。

3. 在“数据驱动测试”窗口中，单击“添加”按钮创建新的数据源。

4. 选择要使用的数据源类型（例如Excel文件、CSV文件、数据库连接）。

5. 选择文件或连接，并根据需要配置设置。

6. 在“数据驱动测试”窗口内选择刚才创建的数据源，并单击“列”按钮将数据映射到测试脚本。

7. 将每个数据列映射到测试脚本中的相应输入参数。

8. 保存更改并关闭“数据驱动测试”窗口。

9. 像往常一样运行您的测试脚本。

10. 在提示时，选择使用多个数据集运行脚本的选项。

11. 脚本将针对数据源中每行数据运行一次，使用该行的值作为输入参数。



#### Part 10 - Verification point

在IBM RFT 8.5中，Verification Point（验证点）是一种用于检查应用程序行为是否符合预期的工具。以下是使用Verification Point进行测试的步骤：

1. 在测试脚本中选择您要验证的对象。

2. 单击菜单栏中的“Insert”（插入）选项，选择“Verification Point”（验证点）。

3. 在弹出窗口中输入验证点名称并选择要验证的属性和方法。

4. 单击“Finish”（完成）后，Verification Point代码将自动添加到脚本中。

5. 运行测试脚本并观察输出结果。

6. 如果Verification Point失败，则会在测试报告中显示，并且可能需要进一步调试问题。

7. 您可以通过更改Verification Point的条件来进行更多测试，或添加其他验证点以确保应用程序的各个方面均得到测试。

请注意，Verification Point不应该是您测试计划中唯一的测试类型。它应该与其他测试技术（如功能测试、性能测试和安全测试）结合使用，以确保应用程序的所有方面均得到充分测试。

#### Part 11 - Verification point with delay

在IBM RFT 8.5中，当您需要等待应用程序执行某些操作后再进行验证时，可以使用Verification Point with Delay（带有延迟的验证点）来做到这一点。以下是使用Verification Point with Delay进行测试的步骤：

1. 在测试脚本中选择要验证的对象。

2. 单击菜单栏中的“Insert”（插入）选项，选择“Verification Point with Delay”（带有延迟的验证点）。

3. 在弹出窗口中输入验证点名称，选择要验证的属性和方法以及延迟时间。

4. 单击“Finish”（完成）后，Verification Point with Delay代码将自动添加到脚本中。

5. 运行测试脚本并观察输出结果。

6. 如果Verification Point with Delay失败，则会在测试报告中显示，并且可能需要进一步调试问题。

与普通的Verification Point相比，Verification Point with Delay允许您设置等待时间，以确保应用程序引起的更改已经生效。这对于包含网络或其他远程操作的应用程序尤其有用，因为这些操作可能需要一些时间才能完成。请注意，在设置延迟时间时，请避免将其设置得过长，以避免增加测试时间。

#### Part 12 - Get text

好的，以下是一份更详细的IBM RFT 8.5中获取应用程序对象文本的步骤说明：

1. 首先，在RFT测试脚本编辑器中选择一个应用程序对象来获取其文本。可以通过手动创建新的测试脚本或打开现有的测试脚本来执行此操作。

2. 在RFT测试脚本编辑器左侧的Object Map面板中，展开“Gui Objects”和应用程序对象所在的层次结构，以便找到要获取文本的对象。

3. 右键单击该对象，并从上下文菜单中选择“Locate”.

4. Object Browser或Test Object Inspector窗口将打开，其中包含选定对象的属性列表。确认该对象可用并具有正确的测试对象类和测试对象属性。

5. 确认位于应用程序对象内部的文本输入区域已具有焦点，并且文本未被选中。这是因为GetText方法将仅返回文本输入区域中当前可见的非选中文本。

6. 输入以下代码以在RFT测试脚本中获取对象文本，并将其存储在变量中：

   ```
   String varName = testObject().getProperty(".text").toString();
   ```

   在这个例子中，“varName”是想要保存文本值的变量名，而“testObject()”是之前从Object Map中选择的应用程序对象。

7. 运行测试脚本并观察输出结果。变量"varName"中将包含选定对象的文本内容。

请注意，GetText方法仅获取当前文本输入区域中可见的非选中文本。如果文本通过滚动或其他方式未在屏幕上完全显示，则可能导致只返回部分文本或不正确的文本。此外，获取对象文本还可能受到应用程序非标准化的字体和格式化方式的影响。因此，在进行字符串比较等文本验证时，请小心使用。

#### Part 13 - Regular expressions
好的，以下是IBM RFT 8.5中使用正则表达式进行GUI测试的具体步骤：

1. 首先，在RFT测试脚本编辑器中选择想要应用正则表达式的对象。可以通过手动创建新的测试脚本或打开现有的测试脚本来执行此操作。

2. 在RFT测试脚本编辑器左侧的Object Map面板中，展开“Gui Objects”和对象所在的层次结构，以找到该对象。

3. 确认该对象可用并具有正确的测试对象类和测试对象属性。

4. 打开测试脚本并输入以下代码以在对象的方法中使用正则表达式：

   ```java
   // 定义正则表达式字符串
   String regex = "reg[ex]+";
   
   // 定义对应于正则表达式的测试对象属性值。
   String testPropertyValue = "regex example";
   
   // 使用matches()方法，将对象属性与正则表达式进行匹配
   boolean textMatches = testPropertyValue.matches(regex);
   
   if (textMatches) {
       System.out.println(testPropertyValue + " matches the regular expression " + regex);
   } else {
       System.out.println(testPropertyValue + " does not match the regular expression " + regex);
   }
   ```

5. 运行测试脚本并观察输出结果。根据正则表达式是否与相应对象属性相匹配，输出应显示文本。

在上面的代码片段中，“regex”是定义的正则表达式字符串，而“testPropertyValue”是想要测试的属性值。在这种情况下，“textMatches”变量将存储一个布尔值，指示对象属性是否与正则表达式匹配。

需要注意的是，在使用正则表达式时，请确保在字符串中进行适当的转义和转义字符处理，以便有效地捕获特殊字符和模式匹配。

RFT提供了许多方法来使用正则表达式进行GUI测试。可以在IBM RFT 8.5帮助文档中查找更详细的信息，并且建议在开发测试脚本之前仔细阅读有关正则表达式的相关内容。

#### Part 14 - Regexp and verification point
好的，以下是IBM RFT 8.5中使用正则表达式进行验证点测试的具体步骤：

1. 首先，在RFT测试脚本编辑器中打开要用于创建验证点的测试脚本。

2. 在测试脚本中找到要创建验证点的对象，并右键单击该对象以打开上下文菜单。

3. 在上下文菜单中选择“Insert Verification Point”选项。

4. 在弹出的“Select Properties to Verify”对话框中，输入可以使用正则表达式进行匹配的属性（如文本属性）的名称。 

5. 然后，在“Verification Point Creation Options”对话框中，选择“Use Regular Expression”复选框，以启用正则表达式选项。

6. 输入想要匹配的文本字符串的正则表达式，并设置其他选项，例如大小写敏感性，也可根据需要更改验证点名称。

7. 单击“OK”，系统将自动在测试脚本中创建一个新的验证点，并将其添加到您已指定的属性。

8. 完成设置后，请运行测试脚本并观察输出结果。如果对象属性的值与指定的正则表达式模式匹配，则证明验证点已通过。

在上面的步骤中，“Insert Verification Point”选项可让您使用IBM RFT的验证点向导来创建一个新的验证点。 此向导提供了一系列创建验证点所需的选项，包括是否使用正则表达式进行文本匹配等选项。

需要注意的是，在使用正则表达式时，请仔细考虑所需的模式，并确保在您的代码中使用正确的转义字符和语法。 此外，还可以通过在IBM RFT帮助菜单中查找“verification point”关键字来查阅有关验证点和测试用例设计的更多信息。

#### Part 15 - Regexp and dynamic text

好的，以下是使用正则表达式处理动态文本的具体步骤：

1. 首先，在IBM RFT测试脚本编辑器中找到需要处理动态文本的对象。 

2. 将对象的属性设置为捕获动态文本的正则表达式。例如，可以使用正则表达式来捕获日期、时间和其他可变文本。

3. 使用RFT API中的“getROProperty”方法来获取动态文本对象的属性值。

4. 将获取的文本字符串与预期的文本进行比较。 如果两个字符串匹配，则将测试结果设为通过。

5. 如果文本不匹配，则会发生错误，测试脚本可据此确定测试是否失败。

下面是一个示例，演示如何使用IBM RFT来处理动态文本。假设有一个对象包含当前日期，并且您想要使用该日期进行后续测试。在这种情况下，可以执行以下步骤：

1. 使用正则表达式捕获包含日期的文本字符串。例如，此处可以使用类似于“\d{2}\/\d{2}\/\d{4}”之类的模式来匹配常见的日期格式。

2. 使用getROProperty方法来检索对象的文本属性。

```java
String datePattern = "\\d{2}\\/\\d{2}\\/\\d{4}"; // 定义正则表达式
  
TestObject dateObject = findObject("<DATE_OBJECT>");  // 查找包含日期的对象
  
String dateText = (String) dateObject.getProperty(".text"); // 使用 getROProperty 方法检索对象属性

Pattern pattern = Pattern.compile(datePattern);   // 编译正则表达式模式
Matcher matcher = pattern.matcher(dateText);     // 在字符串中查找匹配项

if (matcher.find()) {
  String date = matcher.group();   // 找到匹配项并提取日期
  // 在此处定义后续测试的变量或执行其他操作
} else {
  logError("Unable to find date pattern");    // 如果找不到日期模式，则记录错误
}
```

在上述代码示例中，我们首先定义了一个正则表达式模式来匹配日期。 接下来，我们检索包含日期的文本对象，并使用getROProperty方法获取该对象的文本属性。 然后，我们使用Pattern和Matcher类来比较获取的文本字符串与预期的日期模式，如果两者匹配，则将日期存储在变量中，以便进行后续测试。 

需要注意的是，在使用正则表达式处理动态文本时，应尽可能针对可能出现的不同文本格式编写通用的正则表达式模式。 否则，测试结果可能会受到意外的影响，从而导致测试失败。

#### Part 16 - Java basics

好的，以下是一些IBM RFT测试脚本中使用的Java基础知识。

1. 变量

在IBM RFT测试脚本中，可以通过声明变量来存储值并在测试过程中使用它们。Java中声明变量的基本语法如下：

```java
<数据类型> <变量名> = <初始值>;
```

例如，在测试脚本中，您可以声明一个整数变量并将其设置为初始值：

```java
int myNumber = 42;
```

2. 条件语句

条件语句用于根据特定条件执行不同的代码块。 在Java中，常见的条件语句是if-else语句。 if-else语句的基本语法如下：

```java
if (<条件>) {
    // 执行语句块1
}
else {
    // 执行语句块2
}
```

例如，在测试脚本中，您可以使用if-else语句根据对象是否可见执行不同的操作：

```java
TestObject object = findObject("<OBJECT>");
if (object.isVisible()) {
    // 对象可见时执行的语句块
}
else {
    // 对象不可见时执行的语句块
}
```

3. 循环语句

循环语句用于重复执行某些代码块，直到满足特定条件。 在Java中，for和while是最常见的循环语句。 for循环的基本语法如下：

```java
for (<初始化>; <循环条件>; <更新>) {
    // 执行语句块
}
```

例如，在测试脚本中，您可以使用for循环来迭代数组元素并执行某些操作：

```java
int[] myArray = {1, 2, 3, 4, 5};
for (int i=0; i<myArray.length; i++) {
    // 访问当前元素并执行某些操作
}
```

while循环的基本语法如下：

```java
while (<循环条件>) {
    // 执行语句块
}
```

例如，在测试脚本中，您可以使用while循环来重复执行某个操作，直到对象变为可用或直到达到特定时间：

```java
long startTime = System.currentTimeMillis();
while (!object.exists() && System.currentTimeMillis() - startTime < 5000) {
    // 执行某些操作，例如等待500毫秒
}
```

上述示例中，我们使用while循环来检查对象是否存在，并等待最多5秒钟。 如果对象在5秒内不存在，则退出循环。

这是一些简单的Java基础知识，可帮助您编写IBM RFT测试脚本。 随着熟练掌握Java的更高级特性，您可以创建更高效和灵活的测试脚本。

#### Part 17 - Java string and integer
好的，以下是一些IBM RFT测试脚本中使用的Java字符串和整数基础知识。

1. 字符串

在IBM RFT测试脚本中，您可以使用Java字符串来存储和操作文本数据。 Java字符串有许多内置的方法，例如compareTo（）和substring（），可用于执行各种操作。 这里是一些Java字符串的基本用法：

```java
// 定义一个新字符串
String myString = "Hello, world!";

// 通过方法获取字符串长度
int length = myString.length();

// 比较两个字符串
int result = myString.compareTo("Goodbye, world!");

// 获取子字符串
String subString = myString.substring(0, 5);

// 将字符串转换为大写或小写
String upperCaseString = myString.toUpperCase();
String lowerCaseString = myString.toLowerCase();
```

2. 整数

在IBM RFT测试脚本中，您可以使用Java整数来存储和操作数字数据。 Java整数有许多内置的算术和逻辑运算符，例如加号（+）和等于号（==），可用于执行各种操作。 这里是一些Java整数的基本用法：

```java
// 定义一个新整数
int myNumber = 42;

// 执行加法操作
int sum = myNumber + 10;

// 执行减法操作
int difference = myNumber - 10;

// 执行乘法操作
int product = myNumber * 10;

// 执行除法操作
int quotient = myNumber / 10;

// 执行比较操作
boolean isLessThan = myNumber < 50;
boolean isEqualTo = myNumber == 42;
boolean isGreaterThan = myNumber > 30;
```

这是一些简单的Java字符串和整数知识，可帮助您编写IBM RFT测试脚本。 随着熟练掌握这些基础知识和其他高级特性，您可以创建更高效和灵活的测试脚本。

#### Part 18 - Java loop
好的，以下是一些IBM RFT测试脚本中使用的Java循环基础知识。

在编写IBM RFT测试脚本时，您将需要编写代码来执行重复任务。 Java提供了多种类型的循环结构，包括for，while和do-while循环。 这里是每种循环的基本语法：

1. for循环

for循环用于按照指定次数重复执行代码块。 它由一个初始化子句（定义计数器变量），一个条件子句（定义计数器达到目标值时退出循环的条件）和一个增量子句（递增计数器变量）组成。 这里是一个示例for循环：

```java
for (int i = 0; i < 10; i++) {
    System.out.println("The value of i is: " + i);
}
```

这个循环会打印数字0到9。

2. while循环

while循环用于在条件为真时重复执行代码块。 它只有一个条件子句，并在每次循环迭代之前检查该条件。 这里是一个示例while循环：

```java
int i = 0;
while (i < 10) {
    System.out.println("The value of i is: " + i);
    i++;
}
```

这个循环也会打印数字0到9。

3. do-while循环

do-while循环类似于while循环，但它先执行代码块，然后检查条件。 这确保了代码块至少执行一次。 这里是一个示例do-while循环：

```java
int i = 0;
do {
    System.out.println("The value of i is: " + i);
    i++;
} while (i < 10);
```

这个循环也会打印数字0到9。

这些基本的Java循环结构可用于在IBM RFT测试脚本中执行重复任务。 随着熟练掌握这些基础知识和其他高级特性，您可以创建更高效和灵活的测试脚本。

#### Part 19 - Java conditions 
好的，以下是一些IBM RFT测试脚本中使用的Java条件基础知识。

在编写IBM RFT测试脚本时，您将需要编写代码来根据不同条件执行不同的操作。 Java提供了多种类型的条件语句，包括if/else和switch语句。这里是每种语句的基本语法:

1. if/else语句

if/else语句用于基于一个或多个条件执行不同的代码块。它由一个if子句（定义条件），可选的一个或多个else if子句（定义其他条件）以及一个可选的else子句（定义任何其他情况下要执行的代码）组成。 以下是if/else语句的示例代码：

```java
int x = 10;
if (x > 0) {
    System.out.println("x is positive");
} else if (x < 0) {
    System.out.println("x is negative");
} else {
    System.out.println("x is zero");
}
```

2. switch语句

switch语句类似于if/else语句，但它根据一个表达式的值执行不同的代码块。 它由一个switch子句（定义表达式），多个case子句（定义表达式等于某个值时要执行的代码）以及一个可选的default子句（定义任何其他情况下要执行的代码）组成。 以下是switch语句的示例代码：

```java
int dayOfWeek = 4;
switch (dayOfWeek) {
    case 1:
        System.out.println("Monday");
        break;
    case 2:
        System.out.println("Tuesday");
        break;
    case 3:
        System.out.println("Wednesday");
        break;
    case 4:
        System.out.println("Thursday");
        break;
    case 5:
        System.out.println("Friday");
        break;
    default:
        System.out.println("Weekend day");
        break;
}
```

这些基本的Java条件结构可以在IBM RFT测试脚本中执行不同的操作，从而实现更灵活和高效的测试方案。通过熟练掌握这些基础知识和其他高级特性，您可以编写更复杂的测试脚本来验证应用程序是否按预期工作。

#### Part 20 - Apply logic in script
在IBM RFT测试脚本中应用逻辑非常重要。下面是一些有关在IBM RFT脚本中应用逻辑的提示：

1. 使用流程控制语句

使用流程控制语句（例如if/else和for/while循环）来控制代码和执行流程。这些控制结构允许您根据不同情况选择执行不同的代码块，或者多次执行相同的代码块。

以下是一个使用if/else语句的示例，该语句根据条件执行不同的操作：

```java
if (someCondition) {
    doSomething();
} else {
    doSomethingElse();
}
```

以下是一个使用for循环的示例，该循环将重复执行一些操作，直到达到指定条件：

```java
for (int i = 0; i < someLimit; i++) {
    doSomethingRepeatedly();
}
```

2. 将方法分解为更小的部分

将方法分解为更小的部分可以使代码更易于阅读和维护。这也有助于避免代码重复并增加测试脚本的可重用性。如果您发现自己编写了大量的代码行，请考虑将其拆分为更小、更易于管理的部分。

3. 利用断言

利用断言来验证测试结果是否正确。断言是一种在代码中添加的陈述，它表明某个条件必须为真，否则会引发错误并停止代码的执行。这是一种非常有用的技术，可以在脚本运行时快速识别问题。

以下是一个使用断言的示例：

```java
assertTrue(someResult == expectedValue);
```

4. 使用注释

使用注释来记录您的代码。这将使您编写的代码更易于理解和维护。好的注释应该清楚地描述每个方法、变量或语句块的功能，并提供指向相关文档或资源的链接（如果有的话）。

以上是一些有关在IBM RFT测试脚本中应用逻辑的提示。通过正确应用逻辑，您可以编写高质量、易于维护的测试脚本，并确保应用程序按预期工作。

#### Part 21 - Call Methods
在IBM RFT测试脚本中，调用方法是一种重要的技术。它允许您从多个地方调用相同的代码块，使您的脚本更易于维护和扩展。

以下是一些有关在IBM RFT脚本中调用方法的提示：

1. 编写可重用的代码

在编写脚本时，请尽量编写可重用的代码。这意味着编写可以从多个地方调用的通用代码块，而不是针对特定场景编写硬编码解决方案。例如，如果您希望测试应用程序中的搜索功能，则可以编写一个可重用的搜索方法，该方法可以从多个测试用例中调用。

2. 将方法保存在库文件中

为了使您的方法可在整个项目中共享并易于管理，请将其保存在库文件中。库文件只需要编写一次，然后可以在多个测试脚本中使用。这使得您可以更轻松地实现一致性，并减少了代码重复的风险。

3. 使用对象映射

使用对象映射来引用UI元素和控件。以这种方式编写测试脚本，可以使它们更加模块化，并且可以使任何更改都更容易进行。这可以确保您的测试脚本不会因为小的页面布局或控件更改而中断。

以下是一个使用对象映射和方法调用的示例：

```java
// 引用对象映射中的搜索框和搜索按钮
TextGuiSubitemTestObject searchBox = (TextGuiSubitemTestObject)find("searchBox");
GuiTestObject searchButton = find("searchButton");

// 编写可重用的搜索方法
public void search(String keyword) {
    // 在搜索框中输入关键词
    searchBox.setText(keyword);

    // 单击搜索按钮
    searchButton.click();
}

// 从测试脚本中调用搜索方法
search("apple");
```

通过正确调用方法，您将能够在IBM RFT测试脚本中编写更加高效、可维护和灵活的代码。

#### Part 20 - File handing

当您编写代码时，文件处理是一种非常常见的任务。其中一些常见任务包括从文件中读取数据，将数据写入文件以及重命名或删除文件。在Java中，文件处理通常涉及使用File类和必要的流来执行与文件相关的任务。

## 使用File类

在Java中，要处理文件，首先需要创建一个`File`对象，可以通过提供文件路径或文件的封装URI（Uniform Resource Identifier）来实现。例如，下面的示例代码展示了如何在Java中创建一个File对象：

```java
File myFile = new File("path/to/my/file.txt");
```

在这个例子中，我们提供了文件路径，它指向我们想要读取或写入的文件。

## 读取文件

如果要从文件中读取数据，则需要使用Java提供的某种I/O流。有几种不同类型的流，可以根据您要读取的数据类型以及要使用的功能来选择不同的流。

例如，如果要读取文本文件，可以使用`FileReader`对象。请注意，为了确保在读取完文件后关闭文件，我们建议使用try-with-resources语句块来自动关闭文件。下面的例子演示了如何读取文件并将其打印到控制台上：

```java
try (BufferedReader br = new BufferedReader(new FileReader(myFile))) {
    String line;
    while ((line = br.readLine()) != null) {
        System.out.println(line);
    }
} catch (IOException e) {
    System.err.println("Error reading file: " + e.getMessage());
}
```

在这个例子中，我们使用`BufferedReader`来包装一个`FileReader`对象并读取文件的每一行。请注意，在完成读取后，try-with-resources语句块将自动关闭`BufferedReader`和`FileReader`对象。

## 写入文件

如果要将数据写入文件，则需要使用Java提供的另一种I/O流。与读取文件时相同，有多种不同类型的流可用于根据您要写入的数据类型以及要使用的功能进行选择。

例如，如果要写入文本文件，则可以使用`BufferedWriter`对象。下面的代码片段演示了如何打开文件以进行写入：

```java
try (BufferedWriter bw = new BufferedWriter(new FileWriter(myFile))) {
    bw.write("Hello, World!");
} catch (IOException e) {
    System.err.println("Error writing to file: " + e.getMessage());
}
```

在这个例子中，我们使用`BufferedWriter`对象来包装一个`FileWriter`对象。然后，我们调用`write()`方法将字符串写入缓冲区，并通过调用`flush()`方法将缓冲区中的内容刷新到磁盘上的文件中。

## 重命名和删除文件

要重命名或删除文件，可以使用File类中提供的方法来执行这些任务。要重命名文件，请使用`renameTo()`方法，并传递新文件名称作为参数。

例如，下面的代码片段演示了如何将文件从“myfile.txt”重命名为“newfile.txt”：

```java
File myFile = new File("path/to/my/file.txt");
if (myFile.renameTo(new File("path/to/newfile.txt"))) {
    System.out.println("File renamed successfully.");
} else {
    System.err.println("Failed to rename file.");
}
```

要删除文件，请使用`delete()`方法。例如，下面的代码片段演示了如何删除文件：

```java
File myFile = new File("path/to/my/file.txt");
if (myFile.delete()) {
    System.out.println("File deleted successfully.");
} else {
    System.err.println("Failed to delete file.");
}
```

请注意，在某些操作系统中，您可能需要文件系统权限才能重命名或删除文件。此外，如果文件被另一个程序锁定，则可能无法执行这些任务。
0


