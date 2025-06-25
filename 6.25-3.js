var compose = function(functions) {
    
	return function(x) {
        let len = functions.length
        let temp = x
        for(let i = len - 1; i >= 0 ; i--){
            temp = functions[i](temp)
        }
        return temp
    }
};