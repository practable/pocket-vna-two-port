<template>
  
 <div class='container-sm m-2 background-white border rounded'>
    
    <!-- Drag and Drop elements -->
    <drag-and-drop-components id='dragdropcomponents' header="Calibration Standards" :display='standards' :syncPorts="getSyncPorts" @port1change='updatePort1' @port2change='updatePort2'/>


    <!-- Scan and Save logic -->
    <div class='d-flex flex-row justify-content-center form-check-inline'>
        <label class='txt-primary txt-lg me-2'>Port 1: </label>
          <label class='txt-primary txt-lg me-2'>Short</label>
          <input class="form-check-input mt-2 me-2" type="checkbox" value="" id="shortOneCheck" :checked='getShortOneSaved' disabled>
          <label class='txt-primary txt-lg me-2'>Open</label>
          <input class="form-check-input mt-2 me-2" type="checkbox" value="" id="openOneCheck" :checked='getOpenOneSaved' disabled>
          <label class='txt-primary txt-lg me-2'>Load</label>
          <input class="form-check-input mt-2 me-2" type="checkbox" value="" id="loadOneCheck" :checked='getLoadOneSaved' disabled>
          <label class='txt-primary txt-lg me-2'>Thru</label>
          <input class="form-check-input mt-2 me-2" type="checkbox" value="" id="thruCheck" :checked='getThruSaved' disabled>
    </div>

    <div class='d-flex flex-row justify-content-center form-check-inline'>
        <label class='txt-primary txt-lg me-2'>Port 2: </label>
          <label class='txt-primary txt-lg me-2'>Short</label>
          <input class="form-check-input mt-2 me-2" type="checkbox" value="" id="shortTwoCheck" :checked='getShortTwoSaved' disabled>
          <label class='txt-primary txt-lg me-2'>Open</label>
          <input class="form-check-input mt-2 me-2" type="checkbox" value="" id="openTwoCheck" :checked='getOpenTwoSaved' disabled>
          <label class='txt-primary txt-lg me-2'>Load</label>
          <input class="form-check-input mt-2 me-2" type="checkbox" value="" id="loadTwoCheck" :checked='getLoadTwoSaved' disabled>
    </div>
  

    <div class='d-flex flex-row justify-content-center'>
          <button id="scan" type='button' class="button-lg button-primary" @click='scan' :disabled="port1 === '' && port2 === ''">Scan</button>
          <button id="save_to_calibrate" type='button' class="button-lg button-secondary" @click='save' :disabled="!getShowSave">Save</button>
          <button id="request_calibration" type='button' class="button-lg button-tertiary" @click="rangeFreqCalibration" :disabled='!ready_to_calibrate'>Calibrate</button>
    </div>

     



    <transition name='fade'>
      <div v-if='getShowCalibrationModal' class="modal" id='modal-show' tabindex="-1">
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header" style="background-color: #ccc">
              <h5 class="modal-title">Calibrating Pocket VNA</h5>
            </div>
            <div class="modal-body">
              <div class='d-flex row align-items-center'>
                <div class='col-2'>
                  <div class="spinner-border text-primary text-center" role="status">
                    <span class="visually-hidden">Loading...</span>
                  </div>
                </div>
                <div class='col-10'>
                  <p> Please wait for the pocket VNA to calibrate.</p>
                  <p>This could take around 30 seconds.</p>
                </div>
              </div>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-danger" @click="setShowCalibrationModal(false)">Close</button>
            </div>
          </div>
        </div>
      </div>
      </transition>

      <transition name='fade'>
      <div v-if='getShowScanningModal' class="modal" id='modal-show' tabindex="-1">
        <div class="modal-dialog">
          <div class="modal-content">
            <div class="modal-header" style="background-color: #ccc">
              <h5 class="modal-title">Scanning Pocket VNA ports</h5>
            </div>
            <div class="modal-body">
              <div class='d-flex row align-items-center'>
                <div class='col-2'>
                  <div class="spinner-border text-primary text-center" role="status">
                    <span class="visually-hidden">Loading...</span>
                  </div>
                </div>
                <div class='col-10'>
                  <p> Please wait for the pocket VNA to scan ports.</p>
                  <p>This could take around 30 seconds.</p>
                </div>
              </div>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-danger" @click="setShowScanningModal(false)">Close</button>
            </div>
          </div>
        </div>
      </div>
      </transition>

</div>


</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import DragAndDropComponents from '@/components/DragAndDropComponents.vue';
import dayjs from "dayjs";

export default {
   name: 'DragAndDropCalibrate',
    props: ['calibrationState'],
    components:{
        DragAndDropComponents,
    },
    data () {
        return {
            // frequency_points: 20,
            // frequency: 1.0,   //MHz
            // frequency_end: 4.0, //MHz
            // minFrequency: 0,
            // maxFrequency: 4000.0, //MHz
            // units: 1E6, //MHz
            // avgCounts: 1,
            // islog: false,
            // s11: this.sparams.includes('s11') ? true:false,         
            // s12: this.sparams.includes('s12') ? true:false,
            // s21: this.sparams.includes('s21') ? true:false,
            // s22: this.sparams.includes('s22') ? true:false,

            standards: [
                {type: 'short', img: "Short"},
                {type: 'open', img: "Open"},
                {type: 'load', img: "Load"},
                {type: 'thru', img: "Thru"},
            ],
            calibration: [
                {type: 'short', port:"1", required: true, scanned: false, saved: false},
                {type: 'open', port:"1", required: true, scanned: false, saved: false},
                {type: 'load', port:"1", required: true, scanned: false, saved: false},
                {type: 'thru', port:"1", required: true, scanned: false, saved: false},
                {type: 'short', port:"2", required: true, scanned: false, saved: false},
                {type: 'open', port:"2", required: true, scanned: false, saved: false},
                {type: 'load', port:"2", required: true, scanned: false, saved: false}
            ],
            port1: '',
            port2: '',

        }
    },
    computed:{
        ...mapGetters([
            'getCalibrated',
            'getShowCalibrationModal',
            'getSyncPorts',
            'getShowScanningModal'
        ]),
        ready_to_calibrate(){
            let ready = true;
            this.calibration.forEach((connection) => {
                if(connection.required){
                    if(!connection.saved){
                        ready = false;
                    } 
                }
            });
            return ready;
        },
        getShortOneSaved(){
            if(this.calibration[0].saved){
                return true;
            } else{
                return false;
            }
        },
        getShortTwoSaved(){
            if(this.calibration[4].saved){
                return true;
            } else{
                return false;
            }
        },
        getOpenOneSaved(){
            if(this.calibration[1].saved){
                return true;
            } else{
                return false;
            }
        },
        getOpenTwoSaved(){
            if(this.calibration[5].saved){
                return true;
            } else{
                return false;
            }
        },
        getLoadOneSaved(){
            if(this.calibration[2].saved){
                return true;
            } else{
                return false;
            }
        },
        getLoadTwoSaved(){
            if(this.calibration[6].saved){
                return true;
            } else{
                return false;
            }
        },
        getThruSaved(){
            if(this.calibration[3].saved){
                return true;
            } else{
                return false;
            }
        },
        getShowSave(){
            if(this.port1 == 'short' && this.calibration[0].scanned){
                return true;
            } else if(this.port1 == 'open' && this.calibration[1].scanned){
                return true;
            } else if(this.port1 == 'load' && this.calibration[2].scanned){
                return true;
            } else if(this.port2 == 'short' && this.calibration[4].scanned){
                return true;
            } else if(this.port2 == 'open' && this.calibration[5].scanned){
                return true;
            } else if(this.port2 == 'load' && this.calibration[6].scanned){
                return true;
            } else if((this.port1 == 'thru' || this.port2 == 'thru') && this.calibration[3].scanned){
                return true;
            } 
            else {
                return false;
            }
        }
        
    },
    watch:{
        
    },
    created(){
        
    },
    mounted(){

    },
    methods:{
        ...mapActions([
            'setDraggable',
            'setCalibrated',
            'setShowCalibrationModal',
            'setShowScanningModal'
        ]),
        scan(){
            // thru is between port 1 and 2 so if either are connected to thru then scan thru
            if(this.port1 == 'thru' || this.port2 == 'thru'){
                this.calibration[3].scanned = true;
            } else{
                this.calibration.forEach((connection) => {
                    if((connection.port == "1" && connection.type == this.port1) || (connection.port == "2" && connection.type == this.port2)){
                        connection.scanned = true;
                    }
                });
            }

            // ACTUALLY SEND A COMMAND HERE TO RECEIVE CURRENT DATA
            let params = {
                t:dayjs().unix(),
                avg:Number(this.calibrationState.points),
                sparam:{s11:this.calibrationState.sparams.s11,s12:this.calibrationState.sparams.s12,s21:this.calibrationState.sparams.s21,s22:this.calibrationState.sparams.s22},
                what: this.port1 != '' ? this.port1 : this.port2  
            }

            this.$store.dispatch('requestRangeBeforeCal', params);
            this.$store.dispatch('setShowScanningModal', true);
        }, 
        save(){
            // thru is between port 1 and 2 so if either are connected to thru then save thru
            if(this.port1 == 'thru' || this.port2 == 'thru'){
                if(this.calibration[3].scanned){
                    this.calibration[3].saved = true;
                }
            } else {
                this.calibration.forEach((connection) => {
                    if((connection.port == "1" && connection.type == this.port1) || (connection.port == "2" && connection.type == this.port2)){
                        if(connection.scanned){
                            connection.saved = true;
                        }
                    }
                });
            }
        },
        rangeFreqCalibration(){
            //command structure: {"cmd":"rc","range":{"Start":100000,"End":4000000},"size":2,"isLog":true,"avg":1,"sparam":{"S11":true,"S12":false,"S21":true,"S22":false}}
            
            let params = {
                id:'',
                t:dayjs().unix(),
                range:this.calibrationState.range,
                size: Number(this.calibrationState.size),
                islog: this.calibrationState.islog,
                avg:Number(this.calibrationState.avg),
                sparam:this.calibrationState.sparam
            }

            this.$store.dispatch('requestCalibration', params);
            this.$store.dispatch('setShowCalibrationModal', true);

      },
        updatePort1(connected){
            if(connected.type){
                this.port1 = connected.type;
            } else{
                this.port1 = '';
            }
        },
        updatePort2(connected){
            if(connected.type){
                this.port2 = connected.type;
            } else{
                this.port2 = '';
            }
        }


    }
}
</script>

<style>

</style>