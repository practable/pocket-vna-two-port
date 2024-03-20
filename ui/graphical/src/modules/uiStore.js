//Store for variables that are common amongst multiple UI components. 


const uiStore = {
    state: () => ({
        isDraggable: true,
        showParametersSetModal: false,
        showScanningModal: false,
        showCalibrationModal: false,
        showRequestModal: false,
        showVerifiedModal: false,
        showErrorModal: false,
        errorMessage: '',
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
         SET_SHOW_PARAMETERS_SET_MODAL(state, set){
            state.showParametersSetModal = set;
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
        SET_SHOW_VERIFIED_MODAL(state, set){
            state.showVerifiedModal = set;
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
        },
        RESET_CALIBRATION_PORTS(state){
            state.calibrationPorts.forEach(connection => {
                connection.scanned = false;
                connection.saved = false;
            })
        },
        SET_SHOW_ERROR_MODAL(state, set){
            state.showErrorModal = set;
        },
        SET_ERROR_MESSAGE(state, message){
            state.errorMessage = message;
        }     

       },
       actions:{
        setDraggable(context, draggable){
            context.commit('SET_DRAGGABLE', draggable);
        },
        setShowParametersSetModal(context, set){
            context.commit('SET_SHOW_PARAMETERS_SET_MODAL', set);
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
        setShowVerifiedModal(context, set){
            context.commit('SET_SHOW_VERIFIED_MODAL', set);
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
        },
        resetCalibrationPorts(context){
            context.commit('RESET_CALIBRATION_PORTS');
        },
        setShowErrorModal(context, set){
            context.commit('SET_SHOW_ERROR_MODAL', set);
        },
        setErrorMessage(context, message){
            context.commit('SET_ERROR_MESSAGE', message);
        }


       },
       getters:{
        getDraggable(state){
            return state.isDraggable;
        },
        getShowParametersSetModal(state){
            return state.showParametersSetModal;
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
        getShowVerifiedModal(state){
            return state.showVerifiedModal;
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
        },
        getShowErrorModal(state){
            return state.showErrorModal;
        },
        getErrorMessage(state){
            return state.errorMessage;
        }
          
         
       },  
  
  }

  export default uiStore;
