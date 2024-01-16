//Store for sending commands through the dataSocket


const commandStore = {
    state: () => ({
        dataSocket: null,
        isParametersSet: false,
        isCalibrated: false,     //set to false before deploying.
        isVerified: false,      //set to false before deploying
        
        syncPorts: true,        //for calibration in particular, should dragging a standard onto a port then add that standard to both ports - calibration in the hardware does both at once so defaults to true

       }),
       mutations:{
        SET_DATA_SOCKET(state, socket){
            state.dataSocket = socket;
        },
        // REQUEST_SINGLE(state, params){
        //     state.dataSocket.send(JSON.stringify({
        //         "id":params.id,
        //         "t":params.t,
        //         "cmd":"sq",
        //         "freq":params.freq,
        //         "avg":params.avg,
        //         "sparam":{"s11":params.sparam.s11,"s12":params.sparam.s12,"s21":params.sparam.s21,"s22":params.sparam.s22}
        //     }));
        // },
        // REQUEST_RANGE(state, params){
        //     state.dataSocket.send(JSON.stringify({
        //         "id":params.id,
        //         "t":params.t,
        //         "cmd":"rq",
        //         "range":{"start":params.range.start,"end":params.range.end},
        //         "size":params.size,
        //         "islog":params.islog,
        //         "avg":params.avg,
        //         "what": params.what,
        //         "sparam":{"s11":params.sparam.s11,"s12":params.sparam.s12,"s21":params.sparam.s21,"s22":params.sparam.s22}
        //     }));
        // },
        SEND_CALIBRATION_PARAMETERS(state, params){
            console.log(params);
            state.dataSocket.send(JSON.stringify({
                "id":params.t,
                "t":params.t,
                "cmd":"sc",
                "range":{"start":params.range.start,"end":params.range.end},
                "size":params.size,
                "islog":params.islog,
                "avg":params.avg,
            }));
        },
        REQUEST_RANGE_BEFORE_CAL(state, params){
            console.log(params);
            state.dataSocket.send(JSON.stringify({
                //"id": params.what,
                "t": params.t,
                "cmd":"mc",      
                "what": params.what,
                //"avg":params.avg,
                //"sparam":{"S11":params.sparam.s11,"S12":params.sparam.s12,"S21":params.sparam.s21,"S22":params.sparam.s22} 
            }));

        },
        CONFIRM_CAL(state, params){
            console.log(params);
            state.dataSocket.send(JSON.stringify({
                //"id": params.t,
                "t": params.t,
                "cmd":"cc",      
                //"avg":params.avg,
                //"sparam":{"S11":params.sparam.s11,"S12":params.sparam.s12,"S21":params.sparam.s21,"S22":params.sparam.s22} 
            }));

        },
        REQUEST_CALIBRATION(state, params){
            console.log('calibration request sent');
            state.dataSocket.send(JSON.stringify({
                "id":params.id,
                "t":params.t,
                "cmd":"rc",
                "range":{"start":params.range.start,"end":params.range.end},
                "size":params.size,
                "islog":params.islog,
                "avg":params.avg,
                // "sparam":{"s11":params.sparam.s11,"s12":params.sparam.s12,"s21":params.sparam.s21,"s22":params.sparam.s22}
            }));
        },
        REQUEST_RANGE_AFTER_CAL(state, params){
            console.log(params)
            if(state.isCalibrated){
                console.log('range request sent');
                state.dataSocket.send(JSON.stringify({
                    "id": params.what,
                    "t": params.t,
                    "cmd":"crq",
                    "what": params.what,
                    "avg":params.avg,
                    "sparam":{"S11":params.sparam.s11,"S12":params.sparam.s12,"S21":params.sparam.s21,"S22":params.sparam.s22}  //should all be true
                }));
            } else{
                console.log("Error: need to request calibration first");
            }
            
        },
        SET_PARAMETERS_SET(state, set){
            state.isParametersSet = set;
        },
        SET_CALIBRATED(state, set){
            state.isCalibrated = set;
        },
        SET_VERIFIED(state, set){
            state.isVerified = set;
        },
        SET_PORT_OPEN(state){
            state.dataSocket.send(JSON.stringify({
                "set":"port",
                "to":"open"
            }));
        }
            

       },
       actions:{
        setDataSocket(context, socket){
            context.commit("SET_DATA_SOCKET", socket);
        },
        // requestSingle(context, params){
        //     context.commit('REQUEST_SINGLE', params);
        // },
        // requestRange(context, params){
        //     context.commit('REQUEST_RANGE', params);
        // },
        sendCalibrationParameters(context, params){
            context.commit('SEND_CALIBRATION_PARAMETERS', params);
        },
        requestRangeBeforeCal(context, params){
            context.commit('REQUEST_RANGE_BEFORE_CAL', params);
        },
        confirmCal(context, params){
            context.commit('CONFIRM_CAL', params);
        },
        requestCalibration(context, params){
            context.commit('REQUEST_CALIBRATION', params);
        },
        requestRangeAfterCal(context, params){
            context.commit('REQUEST_RANGE_AFTER_CAL', params);
        },
        setParametersSet(context, set){
            context.commit('SET_PARAMETERS_SET', set);
        },
        setCalibrated(context, set){
            context.commit('SET_CALIBRATED', set);
        },
        setVerified(context, set){
            context.commit('SET_VERIFIED', set);
        },
        setPortOpen(context){
            context.commit('SET_PORT_OPEN');
        },
       },
       getters:{
        getDataSocket(state){
            return state.dataSocket;
        },
        getParametersSet(state){
            return state.isParametersSet;
        },
        getCalibrated(state){
            return state.isCalibrated;
        },
        getVerified(state){
            return state.isVerified;
        },
        getSyncPorts(state){
            return state.syncPorts;
        }
          
       },  
  
  }

  export default commandStore;