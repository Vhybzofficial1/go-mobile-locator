export namespace dto {
	
	export class CarrierCreateReq {
	    key: string;
	    province: string;
	    city: string;
	    isp: string;
	
	    static createFrom(source: any = {}) {
	        return new CarrierCreateReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.province = source["province"];
	        this.city = source["city"];
	        this.isp = source["isp"];
	    }
	}
	export class CarrierData {
	    key: string;
	    province: string;
	    city: string;
	    isp: string;
	
	    static createFrom(source: any = {}) {
	        return new CarrierData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.province = source["province"];
	        this.city = source["city"];
	        this.isp = source["isp"];
	    }
	}
	export class CarrierPageData_mobile_locator_internal_dto_CarrierData_ {
	    total: number;
	    list: CarrierData[];
	
	    static createFrom(source: any = {}) {
	        return new CarrierPageData_mobile_locator_internal_dto_CarrierData_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.list = this.convertValues(source["list"], CarrierData);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CarrierUpdateReq {
	    province: string;
	    city: string;
	    isp: string;
	
	    static createFrom(source: any = {}) {
	        return new CarrierUpdateReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.province = source["province"];
	        this.city = source["city"];
	        this.isp = source["isp"];
	    }
	}
	export class Response__mobile_locator_internal_dto_CarrierData_ {
	    code: number;
	    message: string;
	    data?: CarrierData;
	
	    static createFrom(source: any = {}) {
	        return new Response__mobile_locator_internal_dto_CarrierData_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], CarrierData);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response___uint8_ {
	    code: number;
	    message: string;
	    data: number[];
	
	    static createFrom(source: any = {}) {
	        return new Response___uint8_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.message = source["message"];
	        this.data = source["data"];
	    }
	}
	export class Response_interface____ {
	    code: number;
	    message: string;
	    data: any;
	
	    static createFrom(source: any = {}) {
	        return new Response_interface____(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.message = source["message"];
	        this.data = source["data"];
	    }
	}
	export class Response_mobile_locator_internal_dto_CarrierPageData_mobile_locator_internal_dto_CarrierData__ {
	    code: number;
	    message: string;
	    data: CarrierPageData_mobile_locator_internal_dto_CarrierData_;
	
	    static createFrom(source: any = {}) {
	        return new Response_mobile_locator_internal_dto_CarrierPageData_mobile_locator_internal_dto_CarrierData__(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], CarrierPageData_mobile_locator_internal_dto_CarrierData_);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

