# Fix Mod
go mod tidy
go mod vendor
# Get Commit Message 
read -p "Commit Message: " commit_message
# Commit
git add --all
git commit -m "$commit_message"
# Push
git push -u origin main
# Success Message
echo "[*] Successfully Deployed!"