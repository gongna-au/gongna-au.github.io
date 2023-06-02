---
layout: post
title:  面向对象软件工程
subtitle:
tags: [软件工程]
comments: true
---

### 1.项目概述

个性化学习资源推荐系统是一个为用户提供针对性学习资料推荐的系统。本项目采用前后端分离的架构，后端使用Go语言编写，前端使用JavaScript编写，使用MySQL作为数据存储。项目实现了基于内容的推荐算法，根据用户的交互行为为用户推荐相似度较高的学习资源。

### 2.数据库设计

1. 学习资源

   学习资源表（resources）

   ```sql
   CREATE TABLE `resources` (
     `id` int(11) NOT NULL AUTO_INCREMENT,
     `name` varchar(255) NOT NULL,
     `description` text NOT NULL,
     `url`  varchar(255) NOT NULL,
     `tag`  varchar(100) NOT NULL,
     PRIMARY KEY (`id`)
   ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
   
   ```

   这段SQL代码创建了一个名为 "resources" 的表，包含五个列：

   - "name"：一个必填字段，用于存储资源的名称。
   - "description"：一个必填字段，用于存储资源的详细描述。
   - "url"：一个必填字段，用于存储资源的 URL。
   - "tag"：一个必填字段，用于存储资源的标签或分类。

   该表使用 InnoDB 存储引擎和 utf8mb4 字符集，后者支持比标准 utf8 字符集更广泛的 Unicode 字符。

2. 用户行为表（user_actions）：

   该表用于记录用户与学习资源之间的交互行为，包括ID、用户ID、资源ID和行为

   ```sql
   CREATE TABLE `user_actions` (
     `id` int(11) NOT NULL AUTO_INCREMENT,
     `user_id` int(11) NOT NULL,
     `resource_id` int(11) NOT NULL,
     `action_type` varchar(50) NOT NULL,
     PRIMARY KEY (`id`)
   ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
   ```

### 3.推荐算法

项目实现了基于内容的推荐算法，核心思路如下：

1. 获取用户最近交互的资源：遍历给定的用户行为列表，构建一个包含最近访问过的资源ID的集合。

2. 计算资源之间的相似度：遍历所有资源，针对每个资源，提取其特征（如名称、描述、标签等），然后与用户访问过的资源进行比较，计算它们之间的相似度得分。相似度得分的计算方法是：统计两个资源特征的交集大小，并将其除以两者特征集合大小的平方根。这个值越大，则说明两个资源越相似。

3. 推荐资源，并按评分排序：遍历所有资源，对于每个资源，首先排除那些最近被用户访问过的资源；然后，针对每个用户对该资源的交互历史，给出一个初始的评分；最后，乘上前面计算得到的相似度得分，得到最终的评分。根据评分大小将所有资源排序，返回前N个资源作为推荐结果。

   ```go
   func GetRecommendations(userActions []UserActionModel, resources []ResourceModel, numRecs int) []ResourceModel {
   	// 遍历用户与资源之间的交互记录，获取用户最近交互的资源，并保存在字典 recentResources 中。
   	recentResources := make(map[int]bool)
   	for _, action := range userActions {
   		if _, ok := recentResources[action.ResourceID]; !ok {
   			recentResources[action.ResourceID] = true
   		}
   	}
   
   	// 对每个资源计算它们之间的相似度，并保存在字典 simScores 中。
   	simScores := make(map[int]float64)
   	resourceFeatures := make(map[int]map[string]bool) // 资源的特征，例如标签、描述等
   	for _, r := range resources {
   		resourceFeatures[r.ID] = make(map[string]bool)
   		for _, f := range []string{r.Name, r.Description, r.Tag} {
   			resourceFeatures[r.ID][f] = true
   		}
   
   		// 计算相似度
   		simScore := 0.0
   		for id := range recentResources {
   			if id == r.ID {
   				continue
   			}
   			featureCount := 0
   			for f := range resourceFeatures[id] {
   				if resourceFeatures[r.ID][f] {
   					featureCount++
   				}
   			}
   			simScore += float64(featureCount) / math.Sqrt(float64(len(resourceFeatures[id]))*float64(len(resourceFeatures[r.ID])))
   		}
   		simScores[r.ID] = simScore
   	}
   
   	// 对所有资源进行评分，并按照得分排序，返回前 numRecs 个推荐结果。
   	recommendations := make([]ResourceModel, 0)
   	for _, r := range resources {
   		if _, ok := recentResources[r.ID]; ok {
   			continue // 排除最近用户交互过的资源
   		}
   		score := 0.0
   		for _, action := range userActions {
   			if action.ResourceID == r.ID {
   				score += 1.0 // 用户曾经对该资源进行过交互，给予较高评分
   			}
   		}
   		score *= simScores[r.ID] // 加权得分
   		recommendations = append(recommendations, ResourceModel{
   			ID:          r.ID,
   			Name:        r.Name,
   			Description: r.Description,
   			Url:         r.Url,
   			Tag:         r.Tag,
   			Score:       score,
   		})
   	}
   	sort.Slice(recommendations, func(i, j int) bool { return recommendations[i].Score > recommendations[j].Score })
   	if numRecs < len(recommendations) {
   		return recommendations[:numRecs]
   	} else {
   		return recommendations
   	}
   }
   
   ```

我们实现的关于资源特征的处理方式比较简单，只是将名称、描述和标签拼接在一起后以字符串的方式进行比较。如果使用更加复杂的特征提取方法，如人工标注、文本分析或深度学习等，可能会得到更好的推荐效果。

### 4.项目架构和技术栈

本项目采用MVC（Model-View-Controller）架构，主要使用了以下库：

1. spf13/cast：用于处理类型转换，例如将浮点数转换为整数等。
2. go.uber.org/zap：提供了一种高性能的日志库，用于记录项目运行过程中的日志信息。
3. gin-gonic/gin：是一个用Go语言编写的HTTP Web框架，用于处理HTTP请求和路由分发，简化了Web开发过程。
4. golang-jwt/jwt：用于实现JSON Web Token（JWT）的生成和验证，提高了项目的安全性。
5. gorm：是一个优秀的Go语言ORM（Object-Relational Mapping）库，用于简化数据库操作，提高了代码的可读性和可维护性。

以下是各库在项目中的应用及解决的问题：

1. **spf13/cast**：在项目中，遇到要将不同数据类型进行转换的情况。例如，将字符串类型的数字转换为整数类型。于是使用spf13/cast库能简化这些类型转换操作，提高代码的可读性。
2. **go.uber.org/zap**：项目运行过程中，会遇到异常情况或需要记录关键信息。使用go.uber.org/zap库，可以方便地记录日志，帮助快速定位和解决问题。
3. **gin-gonic/gin**：项目中需要处理用户的HTTP请求并根据不同的请求路径进行相应的业务处理。使用gin-gonic/gin库，可以实现高性能的路由分发和请求处理，简化了Web开发过程。
4. **golang-jwt/jwt**：为了保证用户数据的安全性，需要对用户身份进行验证。使用golang-jwt/jwt库，可以实现JSON Web Token（JWT）的生成和验证，提高了项目的安全性。
5. **gorm**：项目中涉及到大量的数据库操作，使用Go语言的原生SQL语句进行操作可能会导致代码冗长、难以维护。使用gorm库，可以简化数据库操作，提高了代码的可读性和可维护性。

### 5.总结

本报告详细我们组个性化学习资源推荐系统的项目概述、数据库设计、推荐算法、项目架构和技术栈。项目采用前后端分离的架构，后端使用Go语言编写，前端使用JavaScript编写，利用MySQL作为数据存储。项目实现了基于内容的推荐算法，根据用户的交互行为为用户推荐相似度较高的学习资源。通过使用多种库，项目实现了高性能、安全性和易维护性。