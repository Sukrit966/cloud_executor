$(document).ready(function() {
    var editor = ace.edit("container");
    editor.setTheme("ace/theme/tomorrow_night");
    editor.getSession().setTabSize(4);
    editor.setOptions({
      fontSize: "20pt"
    });
    editor.getSession().setValue("#include<iostream>\r\n\r\nusing namespace std;\r\rint main(){\r\n    cout<<\"Hello World\";\r\n}")
    editor.getSession().setMode("ace/mode/c_cpp");
    //style="display:flex;justify-content:center;align-content: center;"
    var dialog = document.querySelector('dialog');
    if (! dialog.showModal) {
      dialogPolyfill.registerDialog(dialog);
    }
    dialog.querySelector('.close').addEventListener('click', function() {
      dialog.close();
    });
    $("#play").click(function(){
        var data = editor.getValue();
        console.log(data)
        var t = +new Date;
	var time = t.toString();
        var packageData = {
          "lang":"cpp",
  	      "fdata": data,
	        "fname": time ,
	        "input":"",
	        "ext":".cpp",
	        //"ccommand":"g++ -o /work/cloud_executer/files/" + time + ".out  ./work/cloud_executer/files/" +  + time  + ".cpp",
		"ccommand":"g++  " + time + ".cpp",
		"ecommand":"./a.out"	       
// "ecommand":"./work/cloud_executer/files/" +time + ".out" 
        }
        console.log(JSON.stringify(packageData))
        if(data.trim().length < 1){
          console.log("Empty Text Field");
        }
        var settings = {
          "async": true,
          "url":"http://167.99.203.49:8080/execute",
          "method": "POST",
          "headers": {
              "content-type": "application/json"
           },
           "processData": false,
           "data":JSON.stringify( packageData)
        }
        
        $.ajax(settings).done(function (response) {
          var error = response.error;
          var output = response.output;
          dialog.querySelector('.content').innerHTML = output;
          dialog.showModal();
        });
    })
    /*
    $("#submitBtn").click(function(){
        var selected = $( "input:checked" ).val(); 
        var data = editor.getValue();
        if(data.trim().length < 1){
        	$('#holder').text("Empty Text Field");
        }
        var settings = {
  			"async": true,
  			"crossDomain": true,
  			"url": "https://mockhttp.herokuapp.com/api/v1/create",
  			"method": "POST",
  			"headers": {
				    "content-type": "application/json"
			   },
  			 "processData": false,
  			 "data": data
		    } */
    /*
		$.ajax(settings).done(function (response) {
			var link = response.final_link;
			console.log(response.final_link);
			$("#target").attr('href',link);
			$("#target").text(link);
		});
    }); */ 
});
