package {{.Package}}.constant;

import lark.core.lang.Description;
import lark.core.lang.EnumValuable;
import lark.core.lang.EnumTitlable;
import lark.core.lang.Enums;

@Description("测试类型")
public enum	TestType implements EnumValuable, EnumTitlable {
    /**
     * 好
     */
    GOOD(1, "好"),
    /**
     * 坏
     */
    BAD(2, "坏");

    private int value;
    private String title;

    private TestType(int value, String title) {
        this.value = value;
        this.title = title;
    }

    /**
     * 获取枚举的 int 值,用于数据保存以及序列化
     *
     * @return 枚举的 int 值
     */
    @Override
    public int value() {
        return this.value;
    }

    /**
     * 获取枚举的显示名称
     *
     * @return 枚举的显示名称
     */
    @Override
    public String title() {
        return this.title;
    }

    /**
     * 根据 int 值构建一个枚举对象
     *
     * @param value 需要构建枚举的 int 的值
     * @return 返回相应 value 值的枚举对象
     */
    public static TestType valueOf(int value) {
        return Enums.valueOf(TestType.class, value);
    }
}