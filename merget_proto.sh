#!/bin/bash

# 设置工作目录为脚本所在目录
cd "$(dirname "$0")"

# 定义文件路径
SERVICES_PROTO="proto/services.proto"
TYPES_DIR="proto/types"
OUTPUT_PROTO="Scheduler.proto"

# 创建临时文件
TMP_FILE=$(mktemp)

# 复制services.proto的文件头（包括syntax、package、option）
sed -n '/^syntax/,/^option.*$/p' "$SERVICES_PROTO" > "$TMP_FILE"
echo "" >> "$TMP_FILE"

# 添加services.proto的消息定义
echo "// ===== Service Base Types =====" >> "$TMP_FILE"
sed -n '/^message/,$ p' "$SERVICES_PROTO" >> "$TMP_FILE"
echo "" >> "$TMP_FILE"

# 遍历types目录下的所有proto文件
for file in "$TYPES_DIR"/*.proto; do
    if [ -f "$file" ]; then
        echo "Processing: $file"
        echo "// ===== $(basename "$file") =====" >> "$TMP_FILE"
        # 提取message定义（排除syntax、package、option行）
        sed -n '/^message/,$ p' "$file" >> "$TMP_FILE"
        echo "" >> "$TMP_FILE"
    else
        echo "Warning: No proto files found in $TYPES_DIR"
    fi
done

# 从services.proto提取service定义
echo "// ===== Service Definition =====" >> "$TMP_FILE"
sed -n '/^service/,$ p' "$SERVICES_PROTO" >> "$TMP_FILE"

# 移动临时文件到目标文件
mv "$TMP_FILE" "$OUTPUT_PROTO"

echo "Proto files merged successfully to $OUTPUT_PROTO"