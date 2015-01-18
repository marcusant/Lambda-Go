$(document).ready(function(){
        function code() {
            $("body, body.logo").slideUp(700, function () {
                $("footer, .logo").css("color", "black");
                $(".logo").text("yee");
                $("body")
                    .append('<iframe width="1" height="1" src="//www.youtube.com/embed/sTMWLW6WyZ4?autoplay=1" frameborder="0" allowfullscreen></iframe>')
                    .css("background", 'url("/static/img/home/yee.jpg") center no-repeat fixed')
                    .css("background-color", "black")
                    .slideDown(700);
            })
        }

        var konami = new Konami();
        konami.code = code;
        konami.load();
});