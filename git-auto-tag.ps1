<#
.SYNOPSIS
    Git 自动提交并递增版本号打 Tag 推送到远程。
.DESCRIPTION
    1. 提交所有本地更改到 Git。
    2. 获取最新的版本 Tag（格式 va.b.c 或 path/va.b.c）。
    3. 根据参数 -lv 递增版本号（1=c+1, 2=b+1, 3=a+1）。
    4. 打上新的 Tag 并推送到远程仓库。
.PARAMETER lv
    版本号递增级别：
    - 1 = 修订号 (c) +1（默认）
    - 2 = 次版本号 (b) +1
    - 3 = 主版本号 (a) +1
.PARAMETER subPath
    子包路径，如 "cmd/xf"，如果不提供则为主包
.EXAMPLE
    .\git-auto-tag.ps1 -lv 1                    # 主包递增修订号 (v1.0.0 → v1.0.1)
    .\git-auto-tag.ps1 -lv 2 -subPath "cmd/xf"  # cmd/xf 子包递增次版本号 (cmd/xf/v1.0.0 → cmd/xf/v1.1.0)
#>

param (
    [ValidateSet(1, 2, 3)]
    [int]$lv = 1,  # 默认递增修订号 (c)
    [string]$subPath = ""  # 子包路径，默认为空表示主包
)

# 1. 提交所有本地更改
git add .
$commitMessage = "Auto commit at $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')"
git commit -m $commitMessage

# 2. 构造 Tag 前缀
$tagPrefix = if ($subPath) { "$subPath/" } else { "" }

# 3. 获取最新的版本 Tag
$tagPattern = if ($tagPrefix) { "${tagPrefix}v*.*.*" } else { "v*.*.*" }
$latestTag = git describe --tags --abbrev=0 --match tagPattern 2>$null
Write-Host "最新 $tagPattern 版本 Tag: $latestTag "
if (-not $latestTag) {
    # 如果没有找到匹配的 Tag，尝试从主版本继承或使用默认值
    if ($tagPrefix) {
        # 对于子包，尝试获取主包的最新版本
        $mainLatestTag = git describe --tags --abbrev=0 --match "v*.*.*" 2>$null
        if ($mainLatestTag) {
            # 使用主包版本作为起点
            $latestTag = "${tagPrefix}$($mainLatestTag.TrimStart('v'))"
            Write-Host "未找到子包 $subPath 的版本 Tag，将从主包版本 $mainLatestTag 开始" -ForegroundColor Yellow
        } else {
            $latestTag = "${tagPrefix}0.0.0"
            Write-Host "未找到任何版本 Tag，将从 $latestTag 开始" -ForegroundColor Yellow
        }
    } else {
        $latestTag = "v0.0.0"  # 如果没有 Tag，默认从 v0.0.0 开始
        Write-Host "未找到版本 Tag，将从 $latestTag 开始" -ForegroundColor Yellow
    }
}

# 4. 解析版本号并递增
$versionString = if ($tagPrefix) {
    $latestTag.Substring($tagPrefix.Length).TrimStart('v')
} else {
    $latestTag.TrimStart('v')
}

$versionParts = $versionString.Split('.')
$a = [int]$versionParts[0]
$b = [int]$versionParts[1]
$c = [int]$versionParts[2]

switch ($lv) {
    1 { $c++; Write-Host "递增修订号 (c)" }
    2 { $b++; $c = 0; Write-Host "递增次版本号 (b)" }
    3 { $a++; $b = 0; $c = 0; Write-Host "递增主版本号 (a)" }
}

$newTag = if ($tagPrefix) {
    "${tagPrefix}v$a.$b.$c"
} else {
    "v$a.$b.$c"
}

Write-Host "当前版本: $latestTag → 新版本: $newTag" -ForegroundColor Green

# 5. 打 Tag 并推送到远程
git tag $newTag
git push origin --tags
git push origin  # 推送提交

Write-Host "已提交并推送 Tag: $newTag" -ForegroundColor Cyan
