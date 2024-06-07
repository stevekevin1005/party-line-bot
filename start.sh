# 啟動伺服器
go run main.go

# 確認Go伺服器是否成功啟動
if [ $? -ne 0 ]; then
    echo "Failed to start Go server. Exiting."
    exit 1
fi