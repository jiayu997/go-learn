window.onload = function(){
	function checkIPValid(ip,titile,empty){
		if (! empty){
		    if (ip == ""){
		    	alert(titile + "IP不能为空");
		    	return false;
		    }
		}
		iplist = ip.split(",");
		var rex = /^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/ ;
		for(var i=0;i<iplist.length;i++){
			var count = 0 ;
			// 判断IP是否合法
			if(! rex.test(iplist[i])){
				alert(iplist[i] + "非法IP") ;
				return false ;
			}

			// IP重复判断
			for (j=0;j<iplist.length;j++){
				if (iplist[i]==iplist[j]){
					count ++ ;
				}
			}
			if(count != 1){
				alert(titile+"IP存在重复");
				return false;
			}

		}
		return true;
	}

	function checkFormData(){
		var control = document.getElementById("control") ;
		var re = checkIPValid(control.value,"K8S主节点",false) ;
		if (!re){
			return false;
		}
		var masternum = control.value.split(",");
		if (!(masternum.length == 1 || masternum.length == 4)){
			alert("K8S主节点IP数量错误");
			return false;
		}

		var nfsharbor = document.getElementById("nfsharbor");
		var re = checkIPValid(nfsharbor.value,"NFS/镜像仓库节点",false);
		if (!re){
			return false;
		}
		var nfsharbornum = nfsharbor.value.split(",").length;
		if (nfsharbornum != 1){
			alert("NFS/镜像仓库节点IP只能为一个");
			return false;
		}

		var clusteradmin = document.getElementById("clusteradmin",false);
		var re = checkIPValid(clusteradmin.value,"平台管理节点");
		if (!re){
			return false;
		}

		// 监控告警数据校验
		var monitor_status = document.getElementById("monitor_status");
		var monitor = document.getElementById("monitor");
		if (monitor_status.value == "true"){
			var re = checkIPValid(monitor.value,"监控告警节点",false);
			if (!re){
				return false;
			}
			var monitor_num = monitor.value.split(",").length;
			if (monitor_num != 1){
				alert("集群监控告警节点IP只能为一个");
				return false;
			}
		}

		// 统一日志数校验
		var log_status = document.getElementById("log_status");
		var log = document.getElementById("log");
		if (log_status.value == "true"){
			var re = checkIPValid(log.value,"统一日志节点",false);
			if (!re){
				return false;
			}
			var log_num = log.value.split(",").length;
			if (log_num != 3){
				alert("统一日志节点IP只能为三个");
				return false;
			}
		}
		// 备份工具数据校验
		var backup_status = document.getElementById("backup_status");
		var backup = document.getElementById("backup");
		if (backup_status.value == "true"){
			var re = checkIPValid(backup.value,"备份工具服务端节点",false);
			if (!re){
				return false;
			}
			var backup_num = backup.value.split(",").length;
			if (backup_num != 1){
				alert("备份工具服务端节点IP只能为一个");
				return false;
			}
		}

		// VIP
		var vip_status = document.getElementById("vip_status");
		var vip = document.getElementById("vip");
		if (vip_status.value == "true"){
			var re = checkIPValid(vip.value,"业务高可用",false);
			if (!re){
				return false;
			}
			var vip_num = vip.value.split(",").length;
			if (vip_num != 1){
				alert("业务高可用IP只能为一个");
				return false;
			}
		}
		return true;
	}

	function getFormdata(){
		// 数据处理
		var control = document.getElementById("control").value.split(",");
		var node = []
		// master节点信息
		if (control.length == 1){
			var master = control[0];
			var masterControl = [];
			var k8svip = "";
		}else if(control.length ==4){
			var master = control[0];
			var masterControl = control.slice(1,3);
			var k8svip = control[3];
		}
		var version = document.getElementById("version").value;
		var nfsharbor = document.getElementById("nfsharbor").value;
		if (node.indexOf(nfsharbor)==-1){
			node.push(nfsharbor)
		}
		var clusteradmin = document.getElementById("clusteradmin").value.split(",");
		for(var i=0;i<clusteradmin.length;i++){
			if (clusteradmin[i]== ""){
				continue
			}
			if(node.indexOf(clusteradmin[i])==-1){
				node.push(clusteradmin[i]);
			}
		}
		var monitor_status = document.getElementById("monitor_status").value;
		if (monitor_status == "false"){
			var monitor = {"monitor": "","enable": false};
		}else{
			var tmp = document.getElementById("monitor").value;
			var monitor = {"monitor": tmp,"enable": true};
			if (node.indexOf(tmp)==-1){
				node.push(tmp);
			}
		}

		var log_status = document.getElementById("log_status");
		if (log_status.value =="false"){
			var log = {"log": [],"enable": false};
		}else{
			var tmp = document.getElementById("log").value.split(",")
			var log = {"log": tmp,"enable": true};
			for(var i=0;i<tmp.length;i++){
				if(node.indexOf(tmp[i])==-1){
					node.push(tmp[i]);
				}
			}
		}
		var backup_status = document.getElementById("backup_status")
		if (backup_status.value == "false"){
			var backup = {"backup": "","enable": false}
		}else{
			var backup = {"backup": document.getElementById("backup").value,"enable": true}
		}
		var vip_status = document.getElementById("vip_status")
		if (vip_status.value == "false"){
			var vip = {"vip": "","enable": false}
		}else{
			var vip = {"vip": document.getElementById("vip").value,"enable": true}
		}
		dataForm = {
			"version": version,
			"master": master,
			"node": node,
			"masterControl": masterControl,
			"k8svip": k8svip,
			"nfsharbor": nfsharbor,
			"clusteradmin": clusteradmin,
			"monitor": monitor,
			"log": log,
			"backup": backup,
			"businvip": vip
		}
		return dataForm;
	}

	// 部署前检查
	var check = document.getElementById("check");
	check.onclick = function(){
		var re = checkFormData();
		if (!re){
			return;
		}	
		var dataFrom = getFormdata();
		console.log(dataFrom)
		// 构建ajax发送数据
		$.ajax({
			url: "/api/v1/deploy",
			type: "POST",
			timeout: "5000",
			contentType: "application/json",
			dataType: "json",
			data: JSON.stringify(dataFrom),
			success: function(msg){
				alert(msg);
				document.getElementById("submit").disabled=false;
			},
			error:function(msg){
				alert(msg.responseText);
				document.getElementById("submit").disabled=true;
			}
		})
	}

	// submit提交检查
	var submit = document.getElementById("submit");
	submit.onclick = function(){
		var log = document.getElementById("logoutput");
		log.style.display = "inline-block";
		var loc = window.location;
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

		// 连接时关闭
		// ws.onclose = function(response){
			// console.log(1)
		// }
	}
}