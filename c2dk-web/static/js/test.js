window.onload = function(){
	function addlog(){
		var ws = new WebSocket("ws://"+loc.host+"/api/v1/deployrun");
		// 连接时打开
		ws.onopen = function(){
			console.log("ws open");
			ws.send("hi");
		}

		// 接收到消息时触发
		ws.onmessage = function(response){
			console.log(response);
			var p = document.createElement("p")
			log.appendChild(p)
			p.innerHTML = response.data
			log.scrollTop = log.scrollHeight
		}

		// 关闭
		ws.onclose = function(response){
			console.log(response)
		}
	}

	var logcheck = document.getElementById("logtest");
	logcheck.onclick = function(){
		var loc = window.location;
		var ws = new WebSocket("ws://"+loc.host+"/testlog");
		var data = {
			"name": "jiayu",
			"password": "boysandgirls",
			"sex": "man"
		}
		// 连接时打开
		ws.onopen = function(){
			document.getElementById("logoutput").innerHTML += "<p id='info'>" + "[INFO] " + "服务端连接成功" + "</p>"
			ws.send(JSON.stringify(data))
		}

		// 接收到消息时触发
		ws.onmessage = function(msg){
			document.getElementById("logoutput").innerHTML += msg.data
			console.log(msg)
		}

		// 关闭
		ws.onclose = function(msg){
			document.getElementById("logoutput").innerHTML += "<p id='info'>" + "[INFO] " + "服务端关闭连接" + "</p>"
			console.log(msg)
		}
	}
	    //构建ajax发送数据
	    // $.ajax({
	    	// url: "/test",
	    	// type: "POST",
	    	// timeout: "5000",
	    	// data: {
				// "token": "123"
			// },
	    	// success: function(msg){
				// document.getElementById("logoutput").innerHTML += msg
	    	// },
	    	// error:function(){
	    		// alert("认证失败");
				// return
	    	// }
	    // })
	// }
}