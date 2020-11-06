package {{.Package}}.spring;

import lark.net.rpc.client.ServiceFactory;
import {{.Package}}.iface.TestService;
import org.springframework.boot.autoconfigure.condition.ConditionalOnMissingBean;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Lazy;
import org.springframework.core.Ordered;
import org.springframework.core.annotation.Order;

@Configuration
@Lazy
@Order(Ordered.LOWEST_PRECEDENCE)
public class ServiceAutoConfig {
    private final String SERVER = "{{.Package}}";

    @Bean
    @ConditionalOnMissingBean
    public TestService testService(ServiceFactory serviceFactory) {
        return serviceFactory.get(SERVER, TestService.class);
    }
}