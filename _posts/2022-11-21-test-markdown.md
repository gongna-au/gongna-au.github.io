---
layout: post
title: DUBBO中的设计模式?
subtitle:
tags: [设计模式]
---

# DUBBO 中的设计模式

Dubbo 框架设计概览

## 1.抽象工厂模式

```java
@SPI("javassist")
public interface ProxyFactory {

    // 针对consumer端，创建出代理对象
    @Adaptive({Constants.PROXY_KEY})
    <T> T getProxy(Invoker<T> invoker) throws RpcException;

    // 针对consumer端，创建出代理对象
    @Adaptive({Constants.PROXY_KEY})
    <T> T getProxy(Invoker<T> invoker, boolean generic) throws RpcException;

    // 针对provider端，将服务对象包装成一个Invoker对象
    @Adaptive({Constants.PROXY_KEY})
    <T> Invoker<T> getInvoker(T proxy, Class<T> type, URL url) throws RpcException;
}


```

```java

public abstract class AbstractProxyFactory implements ProxyFactory {
    private static final Class<?>[] INTERNAL_INTERFACES = new Class<?>[]{
        EchoService.class, Destroyable.class
    };

    private static final ErrorTypeAwareLogger logger = LoggerFactory.getErrorTypeAwareLogger(AbstractProxyFactory.class);

    @Override
    public <T> T getProxy(Invoker<T> invoker) throws RpcException {
        return getProxy(invoker, false);
    }

    @Override
    public <T> T getProxy(Invoker<T> invoker, boolean generic) throws RpcException {
        // when compiling with native image, ensure that the order of the interfaces remains unchanged
        LinkedHashSet<Class<?>> interfaces = new LinkedHashSet<>();
        ClassLoader classLoader = getClassLoader(invoker);

        String config = invoker.getUrl().getParameter(INTERFACES);
        if (StringUtils.isNotEmpty(config)) {
            String[] types = COMMA_SPLIT_PATTERN.split(config);
            for (String type : types) {
                try {
                    interfaces.add(ReflectUtils.forName(classLoader, type));
                } catch (Throwable e) {
                    // ignore
                }

            }
        }

        Class<?> realInterfaceClass = null;
        if (generic) {
            try {
                // find the real interface from url
                String realInterface = invoker.getUrl().getParameter(Constants.INTERFACE);
                realInterfaceClass = ReflectUtils.forName(classLoader, realInterface);
                interfaces.add(realInterfaceClass);
            } catch (Throwable e) {
                // ignore
            }

            if (GenericService.class.equals(invoker.getInterface()) || !GenericService.class.isAssignableFrom(invoker.getInterface())) {
                interfaces.add(com.alibaba.dubbo.rpc.service.GenericService.class);
            }
        }

        interfaces.add(invoker.getInterface());
        interfaces.addAll(Arrays.asList(INTERNAL_INTERFACES));

        try {
            return getProxy(invoker, interfaces.toArray(new Class<?>[0]));
        } catch (Throwable t) {
            if (generic) {
                if (realInterfaceClass != null) {
                    interfaces.remove(realInterfaceClass);
                }
                interfaces.remove(invoker.getInterface());

                logger.error(PROXY_UNSUPPORTED_INVOKER, "", "", "Error occur when creating proxy. Invoker is in generic mode. Trying to create proxy without real interface class.", t);
                return getProxy(invoker, interfaces.toArray(new Class<?>[0]));
            } else {
                throw t;
            }
        }
    }

    public abstract <T> T getProxy(Invoker<T> invoker, Class<?>[] types);
}

```

```java
public class JavassistProxyFactory extends AbstractProxyFactory {
    private static final ErrorTypeAwareLogger logger = LoggerFactory.getErrorTypeAwareLogger(JavassistProxyFactory.class);
    private final JdkProxyFactory jdkProxyFactory = new JdkProxyFactory();

    @Override
    @SuppressWarnings("unchecked")
    public <T> T getProxy(Invoker<T> invoker, Class<?>[] interfaces) {
        try {
            return (T) Proxy.getProxy(interfaces).newInstance(new InvokerInvocationHandler(invoker));
        } catch (Throwable fromJavassist) {
            // try fall back to JDK proxy factory
            try {
                T proxy = jdkProxyFactory.getProxy(invoker, interfaces);
                logger.error(PROXY_FAILED_JAVASSIST, "", "", "Failed to generate proxy by Javassist failed. Fallback to use JDK proxy success. " +
                    "Interfaces: " + Arrays.toString(interfaces), fromJavassist);
                return proxy;
            } catch (Throwable fromJdk) {
                logger.error(PROXY_FAILED_JAVASSIST, "", "", "Failed to generate proxy by Javassist failed. Fallback to use JDK proxy is also failed. " +
                    "Interfaces: " + Arrays.toString(interfaces) + " Javassist Error.", fromJavassist);
                logger.error(PROXY_FAILED_JAVASSIST, "", "", "Failed to generate proxy by Javassist failed. Fallback to use JDK proxy is also failed. " +
                    "Interfaces: " + Arrays.toString(interfaces) + " JDK Error.", fromJdk);
                throw fromJavassist;
            }
        }
    }
}
```

```java
/**
 * JdkRpcProxyFactory
 */
public class JdkProxyFactory extends AbstractProxyFactory {

    @Override
    @SuppressWarnings("unchecked")
    public <T> T getProxy(Invoker<T> invoker, Class<?>[] interfaces) {
        return (T) Proxy.newProxyInstance(invoker.getInterface().getClassLoader(), interfaces, new InvokerInvocationHandler(invoker));
    }

    @Override
    public <T> Invoker<T> getInvoker(T proxy, Class<T> type, URL url) {
        return new AbstractProxyInvoker<T>(proxy, type, url) {
            @Override
            protected Object doInvoke(T proxy, String methodName,
                                      Class<?>[] parameterTypes,
                                      Object[] arguments) throws Throwable {
                Method method = proxy.getClass().getMethod(methodName, parameterTypes);
                return method.invoke(proxy, arguments);
            }
        };
    }

}
```

在 RPC 框架中，客户端发送请求和服务端执行请求的过程都是由代理类来完成的。客户端的代理对象叫做 Client Stub，服务端的代理对象叫做 Server Stub。

在 Dubbo 中用了 ProxyFactory 来创建这 2 个相关的对象，有两种实现一种是基于 jdk 动态代理，一种是基于 javaassist

## 2.**适配器模式**

Dubbo 可以支持多个日志框架，每个日志框架的实现都有对应的 Adapter 类，为什么要搞 Adapter 类呢，因为 Dubbo 中日志接口 Logger 用的是自己的，而实现类是引入的。但这些日志实现类的等级和 Dubbo 中定义的日志等级并不完全一致，例如 JdkLogger 中并没有 trace 和 debug 这个等级，所以要用 Adapter 类把 Logger 中的等级对应到实现类中的合适等级。

```java
/**
 * Logger provider
 */
@SPI(scope = ExtensionScope.FRAMEWORK)
public interface LoggerAdapter {

    Logger getLogger(Class<?> key);

    Logger getLogger(String key);

    Level getLevel();

    void setLevel(Level level);

    File getFile();

    void setFile(File file);
}
```

```java
public class Log4jLoggerAdapter implements LoggerAdapter {

    private File file;

    @SuppressWarnings("unchecked")
    public Log4jLoggerAdapter() {
        try {
            org.apache.log4j.Logger logger = LogManager.getRootLogger();
            if (logger != null) {
                Enumeration<Appender> appenders = logger.getAllAppenders();
                if (appenders != null) {
                    while (appenders.hasMoreElements()) {
                        Appender appender = appenders.nextElement();
                        if (appender instanceof FileAppender) {
                            FileAppender fileAppender = (FileAppender) appender;
                            String filename = fileAppender.getFile();
                            file = new File(filename);
                            break;
                        }
                    }
                }
            }
        } catch (Throwable t) {
        }
    }

    private static org.apache.log4j.Level toLog4jLevel(Level level) {
        if (level == Level.ALL) {
            return org.apache.log4j.Level.ALL;
        }
        if (level == Level.TRACE) {
            return org.apache.log4j.Level.TRACE;
        }
        if (level == Level.DEBUG) {
            return org.apache.log4j.Level.DEBUG;
        }
        if (level == Level.INFO) {
            return org.apache.log4j.Level.INFO;
        }
        if (level == Level.WARN) {
            return org.apache.log4j.Level.WARN;
        }
        if (level == Level.ERROR) {
            return org.apache.log4j.Level.ERROR;
        }
        // if (level == Level.OFF)
        return org.apache.log4j.Level.OFF;
    }

    private static Level fromLog4jLevel(org.apache.log4j.Level level) {
        if (level == org.apache.log4j.Level.ALL) {
            return Level.ALL;
        }
        if (level == org.apache.log4j.Level.TRACE) {
            return Level.TRACE;
        }
        if (level == org.apache.log4j.Level.DEBUG) {
            return Level.DEBUG;
        }
        if (level == org.apache.log4j.Level.INFO) {
            return Level.INFO;
        }
        if (level == org.apache.log4j.Level.WARN) {
            return Level.WARN;
        }
        if (level == org.apache.log4j.Level.ERROR) {
            return Level.ERROR;
        }
        // if (level == org.apache.log4j.Level.OFF)
        return Level.OFF;
    }

    @Override
    public Logger getLogger(Class<?> key) {
        return new Log4jLogger(LogManager.getLogger(key));
    }

    @Override
    public Logger getLogger(String key) {
        return new Log4jLogger(LogManager.getLogger(key));
    }

    @Override
    public Level getLevel() {
        return fromLog4jLevel(LogManager.getRootLogger().getLevel());
    }

    @Override
    public void setLevel(Level level) {
        LogManager.getRootLogger().setLevel(toLog4jLevel(level));
    }

    @Override
    public File getFile() {
        return file;
    }

    @Override
    public void setFile(File file) {

    }

}
```

```java
public class JdkLoggerAdapter implements LoggerAdapter {

    private static final String GLOBAL_LOGGER_NAME = "global";

    private File file;

    public JdkLoggerAdapter() {
        try {
            InputStream in = Thread.currentThread().getContextClassLoader().getResourceAsStream("logging.properties");
            if (in != null) {
                LogManager.getLogManager().readConfiguration(in);
            } else {
                System.err.println("No such logging.properties in classpath for jdk logging config!");
            }
        } catch (Throwable t) {
            System.err.println("Failed to load logging.properties in classpath for jdk logging config, cause: " + t.getMessage());
        }
        try {
            Handler[] handlers = java.util.logging.Logger.getLogger(GLOBAL_LOGGER_NAME).getHandlers();
            for (Handler handler : handlers) {
                if (handler instanceof FileHandler) {
                    FileHandler fileHandler = (FileHandler) handler;
                    Field field = fileHandler.getClass().getField("files");
                    File[] files = (File[]) field.get(fileHandler);
                    if (files != null && files.length > 0) {
                        file = files[0];
                    }
                }
            }
        } catch (Throwable t) {
        }
    }

    private static java.util.logging.Level toJdkLevel(Level level) {
        if (level == Level.ALL) {
            return java.util.logging.Level.ALL;
        }
        if (level == Level.TRACE) {
            return java.util.logging.Level.FINER;
        }
        if (level == Level.DEBUG) {
            return java.util.logging.Level.FINE;
        }
        if (level == Level.INFO) {
            return java.util.logging.Level.INFO;
        }
        if (level == Level.WARN) {
            return java.util.logging.Level.WARNING;
        }
        if (level == Level.ERROR) {
            return java.util.logging.Level.SEVERE;
        }
        // if (level == Level.OFF)
        return java.util.logging.Level.OFF;
    }

    private static Level fromJdkLevel(java.util.logging.Level level) {
        if (level == java.util.logging.Level.ALL) {
            return Level.ALL;
        }
        if (level == java.util.logging.Level.FINER) {
            return Level.TRACE;
        }
        if (level == java.util.logging.Level.FINE) {
            return Level.DEBUG;
        }
        if (level == java.util.logging.Level.INFO) {
            return Level.INFO;
        }
        if (level == java.util.logging.Level.WARNING) {
            return Level.WARN;
        }
        if (level == java.util.logging.Level.SEVERE) {
            return Level.ERROR;
        }
        // if (level == java.util.logging.Level.OFF)
        return Level.OFF;
    }

    @Override
    public Logger getLogger(Class<?> key) {
        return new JdkLogger(java.util.logging.Logger.getLogger(key == null ? "" : key.getName()));
    }

    @Override
    public Logger getLogger(String key) {
        return new JdkLogger(java.util.logging.Logger.getLogger(key));
    }

    @Override
    public Level getLevel() {
        return fromJdkLevel(java.util.logging.Logger.getLogger(GLOBAL_LOGGER_NAME).getLevel());
    }

    @Override
    public void setLevel(Level level) {
        java.util.logging.Logger.getLogger(GLOBAL_LOGGER_NAME).setLevel(toJdkLevel(level));
    }

    @Override
    public File getFile() {
        return file;
    }

    @Override
    public void setFile(File file) {

    }

}
```

```java

public class Slf4jLoggerAdapter implements LoggerAdapter {

    private Level level;
    private File file;

    @Override
    public Logger getLogger(String key) {
        return new Slf4jLogger(org.slf4j.LoggerFactory.getLogger(key));
    }

    @Override
    public Logger getLogger(Class<?> key) {
        return new Slf4jLogger(org.slf4j.LoggerFactory.getLogger(key));
    }

    @Override
    public Level getLevel() {
        return level;
    }

    @Override
    public void setLevel(Level level) {
        this.level = level;
    }

    @Override
    public File getFile() {
        return file;
    }

    @Override
    public void setFile(File file) {
        this.file = file;
    }

}
```

```java
public class JclLoggerAdapter implements LoggerAdapter {

    private Level level;
    private File file;

    @Override
    public Logger getLogger(String key) {
        return new JclLogger(LogFactory.getLog(key));
    }

    @Override
    public Logger getLogger(Class<?> key) {
        return new JclLogger(LogFactory.getLog(key));
    }

    @Override
    public Level getLevel() {
        return level;
    }

    @Override
    public void setLevel(Level level) {
        this.level = level;
    }

    @Override
    public File getFile() {
        return file;
    }

    @Override
    public void setFile(File file) {
        this.file = file;
    }

}
```

```java
public class Log4j2LoggerAdapter implements LoggerAdapter {

    private Level level;

    public Log4j2LoggerAdapter() {

    }

    private static org.apache.logging.log4j.Level toLog4j2Level(Level level) {
        if (level == Level.ALL) {
            return org.apache.logging.log4j.Level.ALL;
        }
        if (level == Level.TRACE) {
            return org.apache.logging.log4j.Level.TRACE;
        }
        if (level == Level.DEBUG) {
            return org.apache.logging.log4j.Level.DEBUG;
        }
        if (level == Level.INFO) {
            return org.apache.logging.log4j.Level.INFO;
        }
        if (level == Level.WARN) {
            return org.apache.logging.log4j.Level.WARN;
        }
        if (level == Level.ERROR) {
            return org.apache.logging.log4j.Level.ERROR;
        }
        return org.apache.logging.log4j.Level.OFF;
    }

    private static Level fromLog4j2Level(org.apache.logging.log4j.Level level) {
        if (level == org.apache.logging.log4j.Level.ALL) {
            return Level.ALL;
        }
        if (level == org.apache.logging.log4j.Level.TRACE) {
            return Level.TRACE;
        }
        if (level == org.apache.logging.log4j.Level.DEBUG) {
            return Level.DEBUG;
        }
        if (level == org.apache.logging.log4j.Level.INFO) {
            return Level.INFO;
        }
        if (level == org.apache.logging.log4j.Level.WARN) {
            return Level.WARN;
        }
        if (level == org.apache.logging.log4j.Level.ERROR) {
            return Level.ERROR;
        }
        return Level.OFF;
    }

    @Override
    public Logger getLogger(Class<?> key) {
        return new Log4j2Logger(LogManager.getLogger(key));
    }

    @Override
    public Logger getLogger(String key) {
        return new Log4j2Logger(LogManager.getLogger(key));
    }

    @Override
    public Level getLevel() {
        return level;
    }

    @Override
    public void setLevel(Level level) {
        this.level = level;
    }

    @Override
    public File getFile() {
        return null;
    }

    @Override
    public void setFile(File file) {
    }
}
```

## 3.工厂方法模式

#### Example1:

Dubbo 可以对结果进行缓存，缓存的策略有很多种，一种策略对应一个缓存工厂类

```java
@SPI("lru")
public interface CacheFactory {

    @Adaptive("cache")
    Cache getCache(URL url, Invocation invocation);

}
```

```java

@Deprecated
public abstract class AbstractCacheFactory implements CacheFactory {

    private final ConcurrentMap<String, Cache> caches = new ConcurrentHashMap<String, Cache>();

    @Override
    public Cache getCache(URL url, Invocation invocation) {
        url = url.addParameter(METHOD_KEY, invocation.getMethodName());
        String key = url.toFullString();
        Cache cache = caches.get(key);
        if (cache == null) {
            caches.put(key, createCache(url));
            cache = caches.get(key);
        }
        return cache;
    }

    protected abstract Cache createCache(URL url);

    @Override
    public org.apache.dubbo.cache.Cache getCache(org.apache.dubbo.common.URL url, org.apache.dubbo.rpc.Invocation invocation) {
        return getCache(new URL(url), new Invocation.CompatibleInvocation(invocation));
    }
}
```

```java
public class ExpiringCacheFactory extends AbstractCacheFactory {
    /**
     * Takes url as an method argument and return new instance of cache store implemented by JCache.
     * @param url url of the method
     * @return ExpiringCache instance of cache
     */
    @Override
    protected Cache createCache(URL url) {
        return new ExpiringCache(url);
    }
}
```

```java
public class JCacheFactory extends AbstractCacheFactory {

    /**
     * Takes url as an method argument and return new instance of cache store implemented by JCache.
     * @param url url of the method
     * @return JCache instance of cache
     */
    @Override
    protected Cache createCache(URL url) {
        return new JCache(url);
    }

}
```

```java
public class ThreadLocalCacheFactory extends AbstractCacheFactory {

    /**
     * Takes url as an method argument and return new instance of cache store implemented by ThreadLocalCache.
     * @param url url of the method
     * @return ThreadLocalCache instance of cache
     */
    @Override
    protected Cache createCache(URL url) {
        return new ThreadLocalCache(url);
    }
}
```

```java
public class LruCacheFactory extends AbstractCacheFactory {

    /**
     * Takes url as an method argument and return new instance of cache store implemented by LruCache.
     * @param url url of the method
     * @return ThreadLocalCache instance of cache
     */
    @Override
    protected Cache createCache(URL url) {
        return new LruCache(url);
    }

}
```

```java
public class LfuCacheFactory extends AbstractCacheFactory {

    /**
     * Takes url as an method argument and return new instance of cache store implemented by LfuCache.
     * @param url url of the method
     * @return ThreadLocalCache instance of cache
     */
    @Override
    protected Cache createCache(URL url) {
        return new LfuCache(url);
    }
}
```

#### Example2:

> 注册中心组件的工厂方法模式

##### dubbo 的注册中心组件的产品体系

抽象产品：Registry，AbstractRegistry，FailbackRegistry 等

具体产品：ZookeeperRegistry，NacosRegistry，MulticastRegistry 等

注册中心的核心抽象产品是 Registry 接口。

按照一般的设计思路，复杂的系统中，接口之下是抽象类，这个组件也是一样的思路，提供了一个 AbstractRegistry 抽象类。

对于 dubbo 里面的注册中心，在 AbstractRegistry 抽象类之下，又增加了一层，FailbackRegistry，用于实现注册中心的注册，订阅，查询，通知等方法的重试，实现故障恢复机制。

它里面联合使用了模板方法和工厂方法模式。

##### 注册中心组件的工厂体系

抽象工厂：RegistryFactory，AbstractRegistryFactory

工厂方法（模板方法）：Registry getRegistry(URL url);

具体工厂：ZookeeperRegistryFactory，NacosRegistryFactory，MulticastRegistryFactory 等

##### getRegistry 方法

AbstractRegistryFactory 的注册中心组件的抽象工厂，它里面的 getRegistry 方法，既是工厂方法模式里面的工厂方法，又是模板方法模式里面的模板方法，这个模板方法里面，调用了这个类里面的另外一个抽象方法 createRegistry(), 各个具体工厂，只需要根据自己的情况，实现各自的 createRegistry()方法，即可实现工作方法模式。

```java
/**
 * AbstractRegistryFactory. (SPI, Singleton, ThreadSafe)
 *
 * @see org.apache.dubbo.registry.RegistryFactory
 */
public abstract class AbstractRegistryFactory implements RegistryFactory, ScopeModelAware {

    @Override
    public Registry getRegistry(URL url) {
        if (registryManager == null) {
            throw new IllegalStateException("Unable to fetch RegistryManager from ApplicationModel BeanFactory. " +
                "Please check if `setApplicationModel` has been override.");
        }

        Registry defaultNopRegistry = registryManager.getDefaultNopRegistryIfDestroyed();
        if (null != defaultNopRegistry) {
            return defaultNopRegistry;
        }

        url = URLBuilder.from(url)
            .setPath(RegistryService.class.getName())
            .addParameter(INTERFACE_KEY, RegistryService.class.getName())
            .removeParameter(TIMESTAMP_KEY)
            .removeAttribute(EXPORT_KEY)
            .removeAttribute(REFER_KEY)
            .build();
        String key = createRegistryCacheKey(url);
        Registry registry = null;
        boolean check = url.getParameter(CHECK_KEY, true) && url.getPort() != 0;
        // Lock the registry access process to ensure a single instance of the registry
        registryManager.getRegistryLock().lock();
        try {
            // double check
            // fix https://github.com/apache/dubbo/issues/7265.
            defaultNopRegistry = registryManager.getDefaultNopRegistryIfDestroyed();
            if (null != defaultNopRegistry) {
                return defaultNopRegistry;
            }
            registry = registryManager.getRegistry(key);
            if (registry != null) {
                return registry;
            }
            //create registry by spi/ioc
            registry = createRegistry(url);
            if (check && registry == null) {
                throw new IllegalStateException("Can not create registry " + url);
            }

            if (registry != null) {
                registryManager.putRegistry(key, registry);
            }
        } catch (Exception e) {
            if (check) {
                throw new RuntimeException("Can not create registry " + url, e);
            } else {
                LOGGER.warn("Failed to obtain or create registry ", e);
            }
        } finally {
            // Release the lock
            registryManager.getRegistryLock().unlock();
        }

        return registry;
    }


    protected abstract Registry createRegistry(URL url);

}
```

总结：dubbo 框架里面的注册中心组件，它里面实现了一个非常典型的工厂方法模式，配合模板方法模式，实现了注册中心的扩展。当需要扩展新的注册中心时，只需要增加一个新的具体工厂类，继承抽象产品 AbstractRegistryFactory，实现它里面的 createRegistry()即可。是联合工厂方法模式和模板方法模式。

## 4.装饰者

Dubbo 中网络传输层用到了 Netty，当我们用 Netty 开发时，一般都是写多个 ChannelHandler，然后将这些 ChannelHandler 添加到 ChannelPipeline 上，就是典型的责任链模式

#### Example1:

但是 Dubbo 考虑到有可能替换网络框架组件，所以整个请求发送和请求接收的过程全部用的都是装饰者模式。即只有 NettyServerHandler 实现的接口是 Netty 中的 ChannelHandler，剩下的接口实现的是 Dubbo 中的 ChannelHandler

如下是服务端消息接收会经过的 ChannelHandler

```java
/**
 * NettyServerHandler.
 */
@io.netty.channel.ChannelHandler.Sharable
public class NettyServerHandler extends ChannelDuplexHandler {
    private static final Logger logger = LoggerFactory.getLogger(NettyServerHandler.class);
    /**
     * the cache for alive worker channel.
     * <ip:port, dubbo channel>
     */
    private final Map<String, Channel> channels = new ConcurrentHashMap<>();

    private final URL url;

    private final ChannelHandler handler;

    public NettyServerHandler(URL url, ChannelHandler handler) {
        if (url == null) {
            throw new IllegalArgumentException("url == null");
        }
        if (handler == null) {
            throw new IllegalArgumentException("handler == null");
        }
        this.url = url;
        this.handler = handler;
    }

    public Map<String, Channel> getChannels() {
        return channels;
    }

    @Override
    public void channelActive(ChannelHandlerContext ctx) throws Exception {
        NettyChannel channel = NettyChannel.getOrAddChannel(ctx.channel(), url, handler);
        if (channel != null) {
            channels.put(NetUtils.toAddressString((InetSocketAddress) ctx.channel().remoteAddress()), channel);
        }
        handler.connected(channel);

        if (logger.isInfoEnabled()) {
            logger.info("The connection of " + channel.getRemoteAddress() + " -> " + channel.getLocalAddress() + " is established.");
        }
    }

    @Override
    public void channelInactive(ChannelHandlerContext ctx) throws Exception {
        NettyChannel channel = NettyChannel.getOrAddChannel(ctx.channel(), url, handler);
        try {
            channels.remove(NetUtils.toAddressString((InetSocketAddress) ctx.channel().remoteAddress()));
            handler.disconnected(channel);
        } finally {
            NettyChannel.removeChannel(ctx.channel());
        }

        if (logger.isInfoEnabled()) {
            logger.info("The connection of " + channel.getRemoteAddress() + " -> " + channel.getLocalAddress() + " is disconnected.");
        }
    }

    @Override
    public void channelRead(ChannelHandlerContext ctx, Object msg) throws Exception {
        NettyChannel channel = NettyChannel.getOrAddChannel(ctx.channel(), url, handler);
        handler.received(channel, msg);
        // trigger qos handler
        ctx.fireChannelRead(msg);
    }


    @Override
    public void write(ChannelHandlerContext ctx, Object msg, ChannelPromise promise) throws Exception {
        super.write(ctx, msg, promise);
        NettyChannel channel = NettyChannel.getOrAddChannel(ctx.channel(), url, handler);
        handler.sent(channel, msg);
    }

    @Override
    public void userEventTriggered(ChannelHandlerContext ctx, Object evt) throws Exception {
        // server will close channel when server don't receive any heartbeat from client util timeout.
        if (evt instanceof IdleStateEvent) {
            NettyChannel channel = NettyChannel.getOrAddChannel(ctx.channel(), url, handler);
            try {
                logger.info("IdleStateEvent triggered, close channel " + channel);
                channel.close();
            } finally {
                NettyChannel.removeChannelIfDisconnected(ctx.channel());
            }
        }
        super.userEventTriggered(ctx, evt);
    }

    @Override
    public void exceptionCaught(ChannelHandlerContext ctx, Throwable cause)
        throws Exception {
        NettyChannel channel = NettyChannel.getOrAddChannel(ctx.channel(), url, handler);
        try {
            handler.caught(channel, cause);
        } finally {
            NettyChannel.removeChannelIfDisconnected(ctx.channel());
        }
    }

}
```

```java
public interface ChannelHandler {

    void handlerAdded(ChannelHandlerContext ctx) throws Exception;

    void handlerRemoved(ChannelHandlerContext ctx) throws Exception;

    void exceptionCaught(ChannelHandlerContext ctx, Throwable cause) throws Exception;

    @Inherited
    @Documented
    @Target(ElementType.TYPE)
    @Retention(RetentionPolicy.RUNTIME)
    @interface Sharable {

    }

}

```

#### Example2:

```java
@SPI(value = "dubbo", scope = ExtensionScope.FRAMEWORK)
public interface Protocol {
    /**
     * Get default port when user doesn't config the port.
     *
     * @return default port
     */
    int getDefaultPort();

    @Adaptive
    <T> Exporter<T> export(Invoker<T> invoker) throws RpcException;

    @Adaptive
    <T> Invoker<T> refer(Class<T> type, URL url) throws RpcException;

    void destroy();

    default List<ProtocolServer> getServers() {
        return Collections.emptyList();
    }
}
```

ProtocolFilterWrapper 类是对 Protocol 类的修饰。DUBBO 是通过 SPI 机制实现装饰器模式，我们以 Protocol 接口进行分析，首先分析装饰器类，抽象装饰器核心要点是实现了 Component 并且组合一个 Component 对象。

```java

/**
 * ListenerProtocol
 */
@Activate(order = 100)
public class ProtocolFilterWrapper implements Protocol {

    private final Protocol protocol;

    public ProtocolFilterWrapper(Protocol protocol) {
        if (protocol == null) {
            throw new IllegalArgumentException("protocol == null");
        }
        this.protocol = protocol;
    }

    @Override
    public int getDefaultPort() {
        return protocol.getDefaultPort();
    }

    @Override
    public <T> Exporter<T> export(Invoker<T> invoker) throws RpcException {
        if (UrlUtils.isRegistry(invoker.getUrl())) {
            return protocol.export(invoker);
        }
        FilterChainBuilder builder = getFilterChainBuilder(invoker.getUrl());
        return protocol.export(builder.buildInvokerChain(invoker, SERVICE_FILTER_KEY, CommonConstants.PROVIDER));
    }

    private <T> FilterChainBuilder getFilterChainBuilder(URL url) {
        return ScopeModelUtil.getExtensionLoader(FilterChainBuilder.class, url.getScopeModel()).getDefaultExtension();
    }

    @Override
    public <T> Invoker<T> refer(Class<T> type, URL url) throws RpcException {
        if (UrlUtils.isRegistry(url)) {
            return protocol.refer(type, url);
        }
        FilterChainBuilder builder = getFilterChainBuilder(url);
        return builder.buildInvokerChain(protocol.refer(type, url), REFERENCE_FILTER_KEY, CommonConstants.CONSUMER);
    }

    @Override
    public void destroy() {
        protocol.destroy();
    }

    @Override
    public List<ProtocolServer> getServers() {
        return protocol.getServers();
    }

}
```

在配置文件中配置装饰器

```java
filter=org.apache.dubbo.rpc.protocol.ProtocolFilterWrapper
listener=org.apache.dubbo.rpc.protocol.ProtocolListenerWrapper
```

通过 SPI 机制加载扩展点时会使用装饰器装饰具体构件：

```java
public class ReferenceConfig<T> extends AbstractReferenceConfig {

    private static final Protocol refprotocol = ExtensionLoader.getExtensionLoader(Protocol.class).getAdaptiveExtension();

    private T createProxy(Map<String, String> map) {
        if (isJvmRefer) {
            URL url = new URL(Constants.LOCAL_PROTOCOL, Constants.LOCALHOST_VALUE, 0, interfaceClass.getName()).addParameters(map);
            invoker = refprotocol.refer(interfaceClass, url);
            if (logger.isInfoEnabled()) {
                logger.info("Using injvm service " + interfaceClass.getName());
            }
        }
    }
}
```

最终生成 refprotocol 为如下对象：

```java
ProtocolFilterWrapper(ProtocolListenerWrapper(InjvmProtocol))
```

## 5.责任链模式

#### Dubbo 的调用链

#### Example1:

代理对象（Client Stub 或者 Server Stub）在执行的过程中会执行所有 Filter 的 invoke 方法，但是这个实现方法是对对象不断进行包装，看起来非常像装饰者模式，但是基于方法名和这个 Filter 的功能，我更觉得这个是责任链模式

```java
private static <T> Invoker<T> buildInvokerChain(final Invoker<T> invoker, String key, String group) {
    Invoker<T> last = invoker;
    // 获取自动激活的扩展类
    List<Filter> filters = ExtensionLoader.getExtensionLoader(Filter.class).getActivateExtension(invoker.getUrl(), key, group);
    if (!filters.isEmpty()) {
        for (int i = filters.size() - 1; i >= 0; i--) {
            final Filter filter = filters.get(i);
            final Invoker<T> next = last;
            last = new Invoker<T>() {

                // 省略部分代码

                @Override
                public Result invoke(Invocation invocation) throws RpcException {
                    // filter 不断的套在 Invoker 上，调用invoke方法的时候就会执行filter的invoke方法
                    Result result = filter.invoke(next, invocation);
                    if (result instanceof AsyncRpcResult) {
                        AsyncRpcResult asyncResult = (AsyncRpcResult) result;
                        asyncResult.thenApplyWithContext(r -> filter.onResponse(r, invoker, invocation));
                        return asyncResult;
                    } else {
                        return filter.onResponse(result, invoker, invocation);
                    }
                }

            };
        }
    }
    return last;
}
```

#### Example2:

生产者和消费者最终执行对象都是过滤器链路最后一个节点，整个链路包含多个过滤器进行业务处理。

```text
生产者过滤器链路
EchoFilter > ClassloaderFilter > GenericFilter > ContextFilter > TraceFilter > TimeoutFilter > MonitorFilter > ExceptionFilter > AbstractProxyInvoker

消费者过滤器链路
ConsumerContextFilter > FutureFilter > MonitorFilter > DubboInvoker
```

ProtocolFilterWrapper 作为链路生成核心通过匿名类方式构建过滤器链路，我们以消费者构建过滤器链路为例:

```java
public class ProtocolFilterWrapper implements Protocol {
    private static <T> Invoker<T> buildInvokerChain(final Invoker<T> invoker, String key, String group) {

        // invoker = DubboInvoker
        Invoker<T> last = invoker;

        // 查询符合条件过滤器列表
        List<Filter> filters = ExtensionLoader.getExtensionLoader(Filter.class).getActivateExtension(invoker.getUrl(), key, group);
        if (!filters.isEmpty()) {
            for (int i = filters.size() - 1; i >= 0; i--) {
                final Filter filter = filters.get(i);
                final Invoker<T> next = last;

                // 构造一个简化Invoker
                last = new Invoker<T>() {
                    @Override
                    public Class<T> getInterface() {
                        return invoker.getInterface();
                    }

                    @Override
                    public URL getUrl() {
                        return invoker.getUrl();
                    }

                    @Override
                    public boolean isAvailable() {
                        return invoker.isAvailable();
                    }

                    @Override
                    public Result invoke(Invocation invocation) throws RpcException {
                        // 构造过滤器链路
                        Result result = filter.invoke(next, invocation);
                        if (result instanceof AsyncRpcResult) {
                            AsyncRpcResult asyncResult = (AsyncRpcResult) result;
                            asyncResult.thenApplyWithContext(r -> filter.onResponse(r, invoker, invocation));
                            return asyncResult;
                        } else {
                            return filter.onResponse(result, invoker, invocation);
                        }
                    }

                    @Override
                    public void destroy() {
                        invoker.destroy();
                    }

                    @Override
                    public String toString() {
                        return invoker.toString();
                    }
                };
            }
        }
        return last;
    }

    @Override
    public <T> Invoker<T> refer(Class<T> type, URL url) throws RpcException {
        // RegistryProtocol不构造过滤器链路
        if (Constants.REGISTRY_PROTOCOL.equals(url.getProtocol())) {
            return protocol.refer(type, url);
        }
        Invoker<T> invoker = protocol.refer(type, url);
        return buildInvokerChain(invoker, Constants.REFERENCE_FILTER_KEY, Constants.CONSUMER);
    }
}
```

## 6.建造者模式

```java

@SPI(value = "default", scope = APPLICATION)
public interface FilterChainBuilder {
    /**
     * build consumer/provider filter chain
     */
    <T> Invoker<T> buildInvokerChain(final Invoker<T> invoker, String key, String group);

    <T> ClusterInvoker<T> buildClusterInvokerChain(final ClusterInvoker<T> invoker, String key, String group);


}
```

```java
@Activate
public class DefaultFilterChainBuilder implements FilterChainBuilder {

    /**
     * build consumer/provider filter chain
     */
    @Override
    public <T> Invoker<T> buildInvokerChain(final Invoker<T> originalInvoker, String key, String group) {
        Invoker<T> last = originalInvoker;
        URL url = originalInvoker.getUrl();
        List<ModuleModel> moduleModels = getModuleModelsFromUrl(url);
        List<Filter> filters;
        if (moduleModels != null && moduleModels.size() == 1) {
            filters = ScopeModelUtil.getExtensionLoader(Filter.class, moduleModels.get(0)).getActivateExtension(url, key, group);
        } else if (moduleModels != null && moduleModels.size() > 1) {
            filters = new ArrayList<>();
            List<ExtensionDirector> directors = new ArrayList<>();
            for (ModuleModel moduleModel : moduleModels) {
                List<Filter> tempFilters = ScopeModelUtil.getExtensionLoader(Filter.class, moduleModel).getActivateExtension(url, key, group);
                filters.addAll(tempFilters);
                directors.add(moduleModel.getExtensionDirector());
            }
            filters = sortingAndDeduplication(filters, directors);

        } else {
            filters = ScopeModelUtil.getExtensionLoader(Filter.class, null).getActivateExtension(url, key, group);
        }

        if (!CollectionUtils.isEmpty(filters)) {
            for (int i = filters.size() - 1; i >= 0; i--) {
                final Filter filter = filters.get(i);
                final Invoker<T> next = last;
                last = new CopyOfFilterChainNode<>(originalInvoker, next, filter);
            }
            return new CallbackRegistrationInvoker<>(last, filters);
        }

        return last;
    }

}
```

## 7.模板方法

模板方法模式定义一个操作中的算法骨架，一般使用抽象类定义算法骨架。抽象类同时定义一些抽象方法，这些抽象方法延迟到子类实现，这样子类不仅遵守了算法骨架约定，也实现了自己的算法。既保证了规约也兼顾灵活性。这就是用抽象构建框架，用实现扩展细节。

DUBBO 源码中有一个非常重要的核心概念 Invoker，我的理解是执行器或者说一个可执行对象，能够根据方法的名称、参数得到相应执行结果。

```java
public abstract class AbstractInvoker<T> implements Invoker<T> {

    @Override
    public Result invoke(Invocation inv) throws RpcException {
        RpcInvocation invocation = (RpcInvocation) inv;
        invocation.setInvoker(this);
        if (attachment != null && attachment.size() > 0) {
            invocation.addAttachmentsIfAbsent(attachment);
        }
        Map<String, String> contextAttachments = RpcContext.getContext().getAttachments();
        if (contextAttachments != null && contextAttachments.size() != 0) {
            invocation.addAttachments(contextAttachments);
        }
        if (getUrl().getMethodParameter(invocation.getMethodName(), Constants.ASYNC_KEY, false)) {
            invocation.setAttachment(Constants.ASYNC_KEY, Boolean.TRUE.toString());
        }
        RpcUtils.attachInvocationIdIfAsync(getUrl(), invocation);

        try {
            return doInvoke(invocation);
        } catch (InvocationTargetException e) {
            Throwable te = e.getTargetException();
            if (te == null) {
                return new RpcResult(e);
            } else {
                if (te instanceof RpcException) {
                    ((RpcException) te).setCode(RpcException.BIZ_EXCEPTION);
                }
                return new RpcResult(te);
            }
        } catch (RpcException e) {
            if (e.isBiz()) {
                return new RpcResult(e);
            } else {
                throw e;
            }
        } catch (Throwable e) {
            return new RpcResult(e);
        }
    }

    protected abstract Result doInvoke(Invocation invocation) throws Throwable;
}
```

AbstractInvoker 作为抽象父类定义了 invoke 方法这个方法骨架，并且定义了 doInvoke 抽象方法供子类扩展，例如子类**InjvmInvoker**、**DubboInvoker**各自实现了 doInvoke 方法。

InjvmInvoker 是本地引用，调用时直接从本地暴露生产者容器获取生产者 Exporter 对象即可。

```java
class InjvmInvoker<T> extends AbstractInvoker<T> {

    @Override
    public Result doInvoke(Invocation invocation) throws Throwable {
        Exporter<?> exporter = InjvmProtocol.getExporter(exporterMap, getUrl());
        if (exporter == null) {
            throw new RpcException("Service [" + key + "] not found.");
        }
        RpcContext.getContext().setRemoteAddress(Constants.LOCALHOST_VALUE, 0);
        return exporter.getInvoker().invoke(invocation);
    }
}
```

DubboInvoker 相对复杂一些，需要考虑同步异步调用方式，配置优先级，远程通信等等。

```java
public class DubboInvoker<T> extends AbstractInvoker<T> {

    @Override
    protected Result doInvoke(final Invocation invocation) throws Throwable {
        RpcInvocation inv = (RpcInvocation) invocation;
        final String methodName = RpcUtils.getMethodName(invocation);
        inv.setAttachment(Constants.PATH_KEY, getUrl().getPath());
        inv.setAttachment(Constants.VERSION_KEY, version);
        ExchangeClient currentClient;
        if (clients.length == 1) {
            currentClient = clients[0];
        } else {
            currentClient = clients[index.getAndIncrement() % clients.length];
        }
        try {
            boolean isAsync = RpcUtils.isAsync(getUrl(), invocation);
            boolean isAsyncFuture = RpcUtils.isReturnTypeFuture(inv);
            boolean isOneway = RpcUtils.isOneway(getUrl(), invocation);

            // 超时时间方法级别配置优先级最高
            int timeout = getUrl().getMethodParameter(methodName, Constants.TIMEOUT_KEY, Constants.DEFAULT_TIMEOUT);
            if (isOneway) {
                boolean isSent = getUrl().getMethodParameter(methodName, Constants.SENT_KEY, false);
                currentClient.send(inv, isSent);
                RpcContext.getContext().setFuture(null);
                return new RpcResult();
            } else if (isAsync) {
                ResponseFuture future = currentClient.request(inv, timeout);
                FutureAdapter<Object> futureAdapter = new FutureAdapter<>(future);
                RpcContext.getContext().setFuture(futureAdapter);
                Result result;
                if (isAsyncFuture) {
                    result = new AsyncRpcResult(futureAdapter, futureAdapter.getResultFuture(), false);
                } else {
                    result = new SimpleAsyncRpcResult(futureAdapter, futureAdapter.getResultFuture(), false);
                }
                return result;
            } else {
                RpcContext.getContext().setFuture(null);
                return (Result) currentClient.request(inv, timeout).get();
            }
        } catch (TimeoutException e) {
            throw new RpcException(RpcException.TIMEOUT_EXCEPTION, "Invoke remote method timeout. method: " + invocation.getMethodName() + ", provider: " + getUrl() + ", cause: " + e.getMessage(), e);
        } catch (RemotingException e) {
            throw new RpcException(RpcException.NETWORK_EXCEPTION, "Failed to invoke remote method: " + invocation.getMethodName() + ", provider: " + getUrl() + ", cause: " + e.getMessage(), e);
        }
    }
}
```

## 8.保护性暂停模式

在多线程编程实践中我们肯定会面临线程间数据交互的问题。在处理这类问题时需要使用一些设计模式，从而保证程序的正确性和健壮性。

保护性暂停设计模式就是解决多线程间数据交互问题的一种模式。

(1) 消费者发送请求

顺着这一个链路跟踪代码：消费者发送请求 > 提供者接收请求并执行，并且将运行结果发送给消费者 > 消费者接收结果。

```java
final class HeaderExchangeChannel implements ExchangeChannel {

    @Override
    public ResponseFuture request(Object request, int timeout) throws RemotingException {
        if (closed) {
            throw new RemotingException(this.getLocalAddress(), null, "Failed to send request " + request + ", cause: The channel " + this + " is closed!");
        }
        Request req = new Request();
        req.setVersion(Version.getProtocolVersion());
        req.setTwoWay(true);
        req.setData(request);
        DefaultFuture future = DefaultFuture.newFuture(channel, req, timeout);
        try {
            channel.send(req);
        } catch (RemotingException e) {
            future.cancel();
            throw e;
        }
        return future;
    }
}

class DefaultFuture implements ResponseFuture {

    // FUTURES容器
    private static final Map<Long, DefaultFuture> FUTURES = new ConcurrentHashMap<>();

    private DefaultFuture(Channel channel, Request request, int timeout) {
        this.channel = channel;
        this.request = request;
        // 请求ID
        this.id = request.getId();
        this.timeout = timeout > 0 ? timeout : channel.getUrl().getPositiveParameter(Constants.TIMEOUT_KEY, Constants.DEFAULT_TIMEOUT);
        FUTURES.put(id, this);
        CHANNELS.put(id, channel);
    }
}
```

(2) 提供者接收请求并执行，并且将运行结果发送给消费者

```java
public class HeaderExchangeHandler implements ChannelHandlerDelegate {

    void handleRequest(final ExchangeChannel channel, Request req) throws RemotingException {
        // response与请求ID对应
        Response res = new Response(req.getId(), req.getVersion());
        if (req.isBroken()) {
            Object data = req.getData();
            String msg;
            if (data == null) {
                msg = null;
            } else if (data instanceof Throwable) {
                msg = StringUtils.toString((Throwable) data);
            } else {
                msg = data.toString();
            }
            res.setErrorMessage("Fail to decode request due to: " + msg);
            res.setStatus(Response.BAD_REQUEST);
            channel.send(res);
            return;
        }
        // message = RpcInvocation包含方法名、参数名、参数值等
        Object msg = req.getData();
        try {

            // DubboProtocol.reply执行实际业务方法
            CompletableFuture<Object> future = handler.reply(channel, msg);

            // 如果请求已经完成则发送结果
            if (future.isDone()) {
                res.setStatus(Response.OK);
                res.setResult(future.get());
                channel.send(res);
                return;
            }
        } catch (Throwable e) {
            res.setStatus(Response.SERVICE_ERROR);
            res.setErrorMessage(StringUtils.toString(e));
            channel.send(res);
        }
    }
}
```

(3) 消费者接收结果

```java
class DefaultFuture implements ResponseFuture {
    private final Lock lock = new ReentrantLock();
    private final Condition done = lock.newCondition();

    public static void received(Channel channel, Response response) {
        try {
            // 取出对应的请求对象
            DefaultFuture future = FUTURES.remove(response.getId());
            if (future != null) {
                future.doReceived(response);
            } else {
                logger.warn("The timeout response finally returned at "
                            + (new SimpleDateFormat("yyyy-MM-dd HH:mm:ss.SSS").format(new Date()))
                            + ", response " + response
                            + (channel == null ? "" : ", channel: " + channel.getLocalAddress()
                               + " -> " + channel.getRemoteAddress()));
            }
        } finally {
            CHANNELS.remove(response.getId());
        }
    }


    @Override
    public Object get(int timeout) throws RemotingException {
        if (timeout <= 0) {
            timeout = Constants.DEFAULT_TIMEOUT;
        }
        if (!isDone()) {
            long start = System.currentTimeMillis();
            lock.lock();
            try {
                while (!isDone()) {

                    // 放弃锁并使当前线程阻塞，直到发出信号中断它或者达到超时时间
                    done.await(timeout, TimeUnit.MILLISECONDS);

                    // 阻塞结束后再判断是否完成
                    if (isDone()) {
                        break;
                    }

                    // 阻塞结束后判断是否超时
                    if(System.currentTimeMillis() - start > timeout) {
                        break;
                    }
                }
            } catch (InterruptedException e) {
                throw new RuntimeException(e);
            } finally {
                lock.unlock();
            }
            // response对象仍然为空则抛出超时异常
            if (!isDone()) {
                throw new TimeoutException(sent > 0, channel, getTimeoutMessage(false));
            }
        }
        return returnFromResponse();
    }

    private void doReceived(Response res) {
        lock.lock();
        try {
            // 接收到服务器响应赋值response
            response = res;
            if (done != null) {
                // 唤醒get方法中处于等待的代码块
                done.signal();
            }
        } finally {
            lock.unlock();
        }
        if (callback != null) {
            invokeCallback(callback);
        }
    }
}
```

## 9.单例模式之双重检查锁模式

单例设计模式可以保证在整个应用某个类只能存在一个对象实例，并且这个类只提供一个取得其对象实例方法，通常这个对象创建和销毁比较消耗资源，例如数据库连接对象等等。

在 DUBBO 服务本地暴露时使用了双重检查锁模式判断 exporter 是否已经存在避免重复创建：

```java
public class RegistryProtocol implements Protocol {

    private <T> ExporterChangeableWrapper<T> doLocalExport(final Invoker<T> originInvoker, URL providerUrl) {
        String key = getCacheKey(originInvoker);
        ExporterChangeableWrapper<T> exporter = (ExporterChangeableWrapper<T>) bounds.get(key);
        if (exporter == null) {
            synchronized (bounds) {
                exporter = (ExporterChangeableWrapper<T>) bounds.get(key);
                if (exporter == null) {
                    final Invoker<?> invokerDelegete = new InvokerDelegate<T>(originInvoker, providerUrl);
                    final Exporter<T> strongExporter = (Exporter<T>) protocol.export(invokerDelegete);
                    exporter = new ExporterChangeableWrapper<T>(strongExporter, originInvoker);
                    bounds.put(key, exporter);
                }
            }
        }
        return exporter;
    }
}
```
