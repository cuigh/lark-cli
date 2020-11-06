package {{.Package}}.dto;

import lark.pb.field.FieldType;
import lark.pb.annotation.ProtoField;
import lark.pb.annotation.ProtoMessage;
import {{.Package}}.constant.TestType;

import java.time.LocalDateTime;

public class TestDto {
    /**
     * 请求参数
     */
    @Setter
    @Getter
    @ProtoMessage(description = "请求参数")
    public static class HelloRequest {
        /**
         * ID
         */
        @ProtoField(order = 1, type = FieldType.INT32, required = true, description = "ID")
        private int id;

        /**
         * 类型
         */
        @ProtoField(order = 2, type = FieldType.ENUM, description = "类型")
        private TestType type;

        /**
         * 时间
         */
        @ProtoField(order = 3, type = FieldType.INT64, description = "时间")
        private LocalDateTime time;
    }

    /**
     * 响应结果
     */
    @Setter
    @Getter
    @ProtoMessage(description = "响应结果")
    public static class HelloResponse {
        /**
         * 结果
         */
        @ProtoField(order = 1, type = FieldType.STRING, description = "结果")
        private String result;
    }
}