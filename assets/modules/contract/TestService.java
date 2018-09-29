package {{.Package}}.iface;

import {{.Package}}.dto.TestDto.*;

/**
 * 测试服务
**/
public interface TestService {
	/**
	 * 测试
	**/
	HelloResponse hello(HelloRequest request);
}