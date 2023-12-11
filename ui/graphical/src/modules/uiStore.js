//Store for variables that are common amongst multiple UI components. 


const uiStore = {
    state: () => ({
        isDraggable: true,
        showScanningModal: false,
        showCalibrationModal: false,
        showRequestModal: false,
        sparams: {'s11':true, 's12': true, 's21':true, 's22':true},       //which params are allowed on the UI
        calibration_state: {sparams:{s11:true, s12:true, s21:true, s22: true}, average: 1, points: 20, start: 1, end: 4},        //for reflecting on verification and measurement tabs
        config_json: '',

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
        }
          
         
       },  
  
  }

  export default uiStore;
