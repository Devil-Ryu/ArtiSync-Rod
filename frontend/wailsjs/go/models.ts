export namespace launcher {
	
	export class Launcher {
	    flags: {[key: string]: string[]};
	
	    static createFrom(source: any = {}) {
	        return new Launcher(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.flags = source["flags"];
	    }
	}

}

