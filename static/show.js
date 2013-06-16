var slide = (function (document, window, undefined) {

    var $ = function (id) {
        return document.getElementById(id)
    };

    var $$ = function (elem) {
        var elem = elem || document;
        return elem.getElementsByTagName("section");
    }

    var stop = function (elem, stop) {
        var v = stop ? "paused" : "running";
        elem.style.webkitAnimationPlayState = v;
    }
    
    var starter = function(){
        this.style.opacity = this.style.opacity==1?0:1;
        this.removeEventListener("webkitAnimationStart");
    }

    var stopper = function () {
        this.style.webkitAnimationPlayState = "paused";
        this.style.webkitAnimationName = "";
        this.style.webkitAnimationDirection = "";
        this.removeEventListener("webkitAnimationEnd", stopper);
    }

    var anim = function(elem, out){
    	if(out){
		if(elem.attributes["aout"]){
			return elem.attributes["aout"].value;
		}
		return "fadeOut";
	}

	if(elem.attributes["ain"]){
		return elem.attributes["ain"].value;
	}
	return "fadeIn"
    }
    
    var steps = $$($("show"));
    
    var reset = function(){
        var e;
        for(var i = 0; i < steps.length; i++){
            e = steps[i];
            e.style.webkitAnimationPlayState = "paused";
            e.style.webkitAnimationName = "";
            e.style.webkitAnimationDirection = "";
            e.style.opacity = 0;
            e.style.zIndex = 0;
            e.style.webkitTimingFunction = "ease-in-out";
            
            console.log(e);
        }
    }

    
    var index = 0;

    var elem;

    (function init() {
        for (var i = 0; i < steps.length; i++) {
            elem = steps[i];
            i == 0 ? elem.style.opacity = 1 : elem.style.opacity = 0;
            if(i == 0 &&elem.attributes["sclass"]){
                $("show").classList.add(elem.attributes["sclass"].value);
            }
            stop(elem, true);
        }
    })(document);


    var api = {};

    api.jump = function (step) {
        if (index == step) return;
        if (step >= steps.length) return;
        if (step < 0) return;
        
        reset(index);

        var now = steps[index];
        now.style.opacity = 1;
        now.style.zIndex = 0;
        var next = steps[step];
        next.style.zIndex = 100;

        var leftout = index < step;
        index = step;
        
        
        if(now.attributes["sclass"]){
            $("show").classList.remove(now.attributes["sclass"].value);
        }
        
        if(next.attributes["sclass"]){
            $("show").classList.add(next.attributes["sclass"].value);
        }

        if (leftout) {
            now.style.webkitAnimationName = anim(now, true);
            now.style.webkitAnimationDuration = "1s";
            next.style.webkitAnimationName = anim(next, false);
            next.style.webkitAnimationDuration = "1s";

            now.addEventListener("webkitAnimationStart", starter, false);

            now.addEventListener("webkitAnimationEnd",stopper, false);

            next.addEventListener("webkitAnimationStart", starter, false);

            next.addEventListener("webkitAnimationEnd",stopper, false);


        } else {
            next.style.webkitAnimationName = anim(next, true);
            next.style.webkitAnimationDuration = "1s";
            next.style.webkitAnimationDirection = "reverse";
            now.style.webkitAnimationName = anim(now, false);
            now.style.webkitAnimationDuration = "1s";
            now.style.webkitAnimationDirection = "reverse";

            now.addEventListener("webkitAnimationStart", starter, false);

            now.addEventListener("webkitAnimationEnd",stopper, false);

            next.addEventListener("webkitAnimationStart", starter, false);

            next.addEventListener("webkitAnimationEnd",stopper, false);
        }

        stop(now, false);
        stop(next, false);
    };
    
    api.next = function(){
        this.jump(index +1);
    };
    
    api.prev = function(){
        this.jump(index -1);
    };


    (function addControlls(document) {
        document.body.addEventListener("keydown", function (event) {
            var c = event.keyCode;
            if (c === 37 || c === 39 || c === 32) event.preventDefault();
        });

        document.body.addEventListener("keyup", function (event) {
            var c = event.keyCode;

            /*left*/
            if (c === 37) {
                api.prev();
                event.preventDefault();
                return
            }

            /*Right or Spacebar*/
            if (c === 39 || c === 32) {
                api.next();
                event.preventDefault();
                return
            }
        });
    })(document);

    window.slide = api;

    return api;

}(document, window));
