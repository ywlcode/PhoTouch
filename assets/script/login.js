$(document).ready(function(){
    var username = $(".name");
    var userpass = $(".pass");
    $(".up").click(function(){
        var pname = username.val();
        var password = userpass.val();
        $.post("/user/login",{name:pname,pwd:password},function(data){
            if(data.ss == '200')
            {
                $(".user *").css("display","none");
                $(".user").css("display","none");
                $(".jump").css("display","block");
                var t=setTimeout(function(){
                    window.location.href='/';
                },800);
            }
            else
            {
                alert("邮箱或密码错误,请重新输入")
            }
        })
    });
})