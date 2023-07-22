window.onload = function(){
	var submit = document.getElementById("submit");
	submit.onclick = function(){
		var username = document.getElementById("username");
		var token = document.getElementById("token");
		if (username.value == ""){
			alert("用户名不能为空");
			return ;
		}
		if (token.value == ""){
			alert("token 不能为空")
			return ;
		}
	    // 构建ajax发送数据
	    $.ajax({
	    	url: "/api/v1/login",
	    	type: "POST",
	    	timeout: "5000",
	    	data: {
				"username": username.value,
				"token": token.value,
			},
	    	success: function(msg){
	    		alert("认证成功");
				window.location.replace("/deploy");
	    	},
	    	error:function(msg){
	    		alert("认证失败");
	    	}
	    })
	}
}