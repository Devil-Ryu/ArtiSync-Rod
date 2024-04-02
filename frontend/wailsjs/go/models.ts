export namespace proto {
	
	export class NetworkCookie {
	    name: string;
	    value: string;
	    domain: string;
	    path: string;
	    expires: number;
	    size: number;
	    httpOnly: boolean;
	    secure: boolean;
	    session: boolean;
	    sameSite?: string;
	    priority: string;
	    sameParty: boolean;
	    sourceScheme: string;
	    sourcePort: number;
	    partitionKey?: string;
	    partitionKeyOpaque?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new NetworkCookie(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.value = source["value"];
	        this.domain = source["domain"];
	        this.path = source["path"];
	        this.expires = source["expires"];
	        this.size = source["size"];
	        this.httpOnly = source["httpOnly"];
	        this.secure = source["secure"];
	        this.session = source["session"];
	        this.sameSite = source["sameSite"];
	        this.priority = source["priority"];
	        this.sameParty = source["sameParty"];
	        this.sourceScheme = source["sourceScheme"];
	        this.sourcePort = source["sourcePort"];
	        this.partitionKey = source["partitionKey"];
	        this.partitionKeyOpaque = source["partitionKeyOpaque"];
	    }
	}
	export class NetworkCookieParam {
	    name: string;
	    value: string;
	    url?: string;
	    domain?: string;
	    path?: string;
	    secure?: boolean;
	    httpOnly?: boolean;
	    sameSite?: string;
	    expires?: number;
	    priority?: string;
	    sameParty?: boolean;
	    sourceScheme?: string;
	    sourcePort?: number;
	    partitionKey?: string;
	
	    static createFrom(source: any = {}) {
	        return new NetworkCookieParam(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.value = source["value"];
	        this.url = source["url"];
	        this.domain = source["domain"];
	        this.path = source["path"];
	        this.secure = source["secure"];
	        this.httpOnly = source["httpOnly"];
	        this.sameSite = source["sameSite"];
	        this.expires = source["expires"];
	        this.priority = source["priority"];
	        this.sameParty = source["sameParty"];
	        this.sourceScheme = source["sourceScheme"];
	        this.sourcePort = source["sourcePort"];
	        this.partitionKey = source["partitionKey"];
	    }
	}

}

