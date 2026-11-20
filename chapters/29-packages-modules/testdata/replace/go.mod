module example.com/shop

go 1.27

require example.com/price v0.0.0

// example.com/price 沒有發佈過；replace 把它換成本機的目錄。
replace example.com/price => ../workspace/price
