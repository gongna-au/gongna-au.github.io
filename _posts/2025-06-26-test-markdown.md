---
layout: post
title: Part1-POM（Project Object Model）
subtitle: 
tags: [java]
comments: true
---

在Java生态中，POM（Project Object Model）文件是Apache Maven构建工具的核心配置文件，通常命名为pom.xml。它通过XML格式定义项目的结构、依赖关系、构建配置和项目元数据。以下是理解POM文件的关键点：

# 1.POM文件的核心作用

- 依赖管理：声明项目所需的第三方库（如JUnit、Spring等），Maven自动从仓库下载并管理传递性依赖。
- 构建生命周期：定义编译、测试、打包、部署等构建流程。
- 项目信息：描述项目的基本信息（如名称、版本、开发者、许可证等）。
- 插件配置：集成Maven插件（如编译器插件、打包插件）并配置其行为。
- 多模块管理：支持聚合多个子模块项目（如微服务架构）。

# 2.POM文件的核心结构
```xml
<project>
    <!-- 基本信息 -->
    <modelVersion>4.0.0</modelVersion>  <!-- 固定值 -->
    <groupId>com.example</groupId>      <!-- 组织标识（反向域名） -->
    <artifactId>my-app</artifactId>     <!-- 项目唯一标识 -->
    <version>1.0.0</version>            <!-- 项目版本 -->
    <packaging>jar</packaging>          <!-- 打包类型（jar/war/pom） -->

    <!-- 依赖管理 -->
    <dependencies>
        <dependency>
            <groupId>junit</groupId>
            <artifactId>junit</artifactId>
            <version>4.12</version>
            <scope>test</scope>         <!-- 依赖作用域 -->
        </dependency>
    </dependencies>

    <!-- 构建配置 -->
    <build>
        <plugins>
            <plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-compiler-plugin</artifactId>
                <version>3.8.1</version>
                <configuration>
                    <source>11</source>  <!-- 指定JDK版本 -->
                    <target>11</target>
                </configuration>
            </plugin>
        </plugins>
    </build>

    <!-- 项目信息 -->
    <name>My Application</name>
    <description>A sample project</description>
</project>
```

# 3.关键元素详解

## 坐标（Coordinates）
groupId + artifactId + version：构成项目的唯一标识（GAV坐标），类似于坐标定位

```xml
<groupId>org.springframework</groupId>
<artifactId>spring-core</artifactId>
<version>5.3.10</version>
```

## 依赖（Dependencies）

- 依赖作用域（Scope）：
  - compile（默认）：编译、测试、运行时均有效。
  - test：仅测试阶段有效（如JUnit）。
  - provided：编译和测试有效，运行时由环境提供（如Servlet API）。
  - runtime：运行时需要，编译时不需要（如JDBC驱动）。
构建配置（Build）

```xml
<plugins>
    <!-- 指定JDK版本 -->
    <plugin>
        <groupId>org.apache.maven.plugins</groupId>
        <artifactId>maven-compiler-plugin</artifactId>
        <configuration>
            <source>11</source>
            <target>11</target>
        </configuration>
    </plugin>
    
    <!-- 打包可执行JAR -->
    <plugin>
        <groupId>org.apache.maven.plugins</groupId>
        <artifactId>maven-shade-plugin</artifactId>
        <version>3.2.4</version>
        <executions>
            <execution>
                <phase>package</phase>
                <goals><goal>shade</goal></goals>
            </execution>
        </executions>
    </plugin>
</plugins>
```

## 属性（Properties）
```xml
<properties>
    <java.version>11</java.version>
    <junit.version>5.7.0</junit.version>
</properties>

<!-- 使用变量 -->
<dependency>
    <groupId>org.junit.jupiter</groupId>
    <artifactId>junit-jupiter</artifactId>
    <version>${junit.version}</version>
</dependency>
```

# 4. 高级特性

## 继承（Inheritance）

- 子模块继承父POM的配置（减少重复）：

```xml
<!-- 子模块pom.xml -->
<parent>
    <groupId>com.example</groupId>
    <artifactId>parent-project</artifactId>
    <version>1.0.0</version>
</parent>
```

## 聚合（Aggregation）

```xml
<!-- 父POM中定义子模块 -->
<modules>
    <module>module-1</module>
    <module>module-2</module>
</modules>
```

## 依赖管理（Dependency Management）

- 父POM统一管理依赖版本，子模块无需重复指定：

```xml
<!-- 父POM中 -->
<dependencyManagement>
    <dependencies>
        <dependency>
            <groupId>com.google.guava</groupId>
            <artifactId>guava</artifactId>
            <version>30.1.1-jre</version>
        </dependency>
    </dependencies>
</dependencyManagement>

<!-- 子模块中直接引用（不写version） -->
<dependency>
    <groupId>com.google.guava</groupId>
    <artifactId>guava</artifactId>
</dependency>
```

# 5.工作流程

1. 编写代码：在src/main/java目录下开发。
2. 声明依赖：在<dependencies>中添加所需库。
3. 构建项目：执行Maven命令：

```xml
mvn clean package   # 清理并打包
mvn test            # 运行测试
mvn install         # 安装到本地仓库

```

# 6.常见问题解决

- 依赖冲突：使用mvn dependency:tree查看依赖树，通过<exclusions>排除冲突库。
- 构建失败：检查JDK版本、网络连接（仓库访问）、插件兼容性。
- 多环境配置：结合profiles实现不同环境（dev/test/prod）的切换。