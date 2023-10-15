cd frontend
npm install
npm run generate

# 確認Nuxt打包是否成功
if [ $? -ne 0 ]; then
    echo "Nuxt build failed. Exiting."
    exit 1
fi

# 切換到Go伺服器目錄並啟動伺服器
cd ../
go run main.go

# 確認Go伺服器是否成功啟動
if [ $? -ne 0 ]; then
    echo "Failed to start Go server. Exiting."
    exit 1
fi