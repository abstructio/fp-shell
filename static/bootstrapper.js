(function(document, json, undefined){
	
	var $ = function(id){
		return document.getElementById(id);
	};

	var $$ = function(elem,attrname, value){
		elem.setAttribute(attrname, value);
	}

	var show = $("show");

	var pres = JSON.parse(json);	

	var init = function(){
		var slide, section, h1;

		for(var i = 0; i < pres.Slides.length; i++){
			slide = pres.Slides[i];

			section = document.createElement("section");

			if(slide.Sclass && slide.Sclass.length != 0){
				$$(section, "sclass", slide.Sclass.join(","));
			}

			if(slide.Class && slide.Class.length != 0){
				$$(section, "class", slide.Class.join(","));
			}

			if(slide.Title && slide.Title != ""){
				h1 = document.createElement("h1");
				h1.innerHTML = slide.Title;
				section.appendChild(h1);
			}

			if(slide.Content && slide.Content != ""){
				section.innerHTML += slide.Content;
			}

			if(slide.Ain && slide.Ain != ""){
				$$(section, "ain", slide.Ain);
			}

			if(slide.Aout && slide.Aout != ""){
				$$(section, "aout", slide.Aout);
			}

			show.appendChild(section);
		}
	}

	init();


})(document, json);
