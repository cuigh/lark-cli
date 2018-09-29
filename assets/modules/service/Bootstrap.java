package {{.Package}};

import lark.net.rpc.RpcApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class Bootstrap {
    public static void main(String[] args) {
        RpcApplication app = new RpcApplication(Bootstrap.class);
        app.run(args);
    }
}