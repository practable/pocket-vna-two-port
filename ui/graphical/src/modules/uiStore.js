//Store for variables that are common amongst multiple UI components. 


const uiStore = {
    state: () => ({
        isDraggable: true,
        showScanningModal: false,
        showCalibrationModal: false,
        showRequestModal: false,
        sparams: {'s11':true, 's12': true, 's21':true, 's22':true},       //which params are allowed on the UI
        calibration_state: {sparam:{s11:true, s12:true, s21:true, s22: true}, islog:false, avg: 1, size: 20, range:{start: 1000000, end: 4000000}},        //set by SetParameters.vue and then read and displayed in other cal, verify and measure tabs
        config_json: '',
        calibrationPorts: [
            {type: 'short', port:"1", required: true, scanned: false, saved: false},
            {type: 'open', port:"1", required: true, scanned: false, saved: false},
            {type: 'load', port:"1", required: true, scanned: false, saved: false},
            {type: 'thru', port:"1", required: true, scanned: false, saved: false},
            {type: 'short', port:"2", required: true, scanned: false, saved: false},
            {type: 'open', port:"2", required: true, scanned: false, saved: false},
            {type: 'load', port:"2", required: true, scanned: false, saved: false}
        ],

       }),
       mutations:{
        SET_DRAGGABLE(state, draggable){
            state.isDraggable = draggable;
         },
         SET_SHOW_SCANNING_MODAL(state, set){
            state.showScanningModal = set;
        },
         SET_SHOW_CALIBRATION_MODAL(state, set){
             state.showCalibrationModal = set;
         },
         SET_SHOW_REQUEST_MODAL(state, set){
            state.showRequestModal = set;
        },
        SET_CONFIG_JSON(state, json){
            state.config_json = json;
        },
        SET_CALIBRATION_STATE(state, params){
            state.calibration_state = params;
        },
        SET_DUT(state, dut){
            state.dut_selected = dut;
        },
        SET_SCANNED(state, what){
            state.calibrationPorts.forEach((connection) => {
                if(connection.type == what){
                    connection.scanned = true;
                }
            });
        }     

       },
       actions:{
        setDraggable(context, draggable){
            context.commit('SET_DRAGGABLE', draggable);
        },
        setShowScanningModal(context, set){
            context.commit('SET_SHOW_SCANNING_MODAL', set);
        },
        setShowCalibrationModal(context, set){
            context.commit('SET_SHOW_CALIBRATION_MODAL', set);
        },
        setShowRequestModal(context, set){
            context.commit('SET_SHOW_REQUEST_MODAL', set);
        },
        setConfigJSON(context, json){
            context.commit('SET_CONFIG_JSON', json);
        },
        setCalibrationState(context, params){
            context.commit('SET_CALIBRATION_STATE', params);
        },
        setDUT(context, dut){
            context.commit('SET_DUT', dut);
        },
        setScanned(context, what){
            context.commit('SET_SCANNED', what);
        }


       },
       getters:{
        getDraggable(state){
            return state.isDraggable;
        },
        getShowScanningModal(state){
            return state.showScanningModal;
        },
        getShowCalibrationModal(state){
            return state.showCalibrationModal;
        },
        getShowRequestModal(state){
            return state.showRequestModal;
        },
        getConfigJSON(state){
            return state.config_json;
        },
        getSParams(state){
            let p_list = [];
            Object.keys(state.sparams).forEach(param => {
                if(state.sparams[param]){
                    p_list.push(param)
                }
            });

            return p_list;
        },
        getCalibrationState(state){
            return state.calibration_state;
        },
        getDUT(state){
            return state.dut_selected;
        },
        getCalibrationPorts(state){
            return state.calibrationPorts;
        }
          
         
       },  
  
  }

  export default uiStore;
