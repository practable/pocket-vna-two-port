<template>
  
 <div class='container-sm m-2 background-white border rounded'>
    <h4 v-if="sparams == 's11'"> 1 port Calibration </h4>
    <h4 v-else> 2 port Calibration </h4>

    <div v-if='sparams.length > 1' class='row mb-2'>
      <div class='col-12'>
        <div v-if='sparams.includes("s11")' class="form-check form-check-inline">
          <input class="form-check-input" type="checkbox" name="flexCheckDefault" id="s11check" checked disabled>
          <label class="form-check-label" for="s11check">S11</label>
        </div>

        <div v-if='sparams.includes("s12")' class="form-check form-check-inline">
          <input class="form-check-input" type="checkbox" name="flexCheckDefault" id="s12check" checked disabled>
          <label class="form-check-label" for="s12check">S12</label>
        </div>

        <div v-if='sparams.includes("s21")' class="form-check form-check-inline">
          <input class="form-check-input" type="checkbox" name="flexCheckDefault" id="s21check" checked disabled>
          <label class="form-check-label" for="s21check">S21</label>
        </div>

         <div v-if='sparams.includes("s22")' class="form-check form-check-inline">
          <input class="form-check-input" type="checkbox" name="flexCheckDefault" id="s22check" checked disabled>
          <label class="form-check-label" for="s22check">S22</label>
        </div>
        
      </div>
    </div>

    <div class='row m-3' @mousedown="setDraggable(false)" @mouseup="setDraggable(true)">
      
        <div class='input-group'>
          <div class='col-md-2 pe-2'>
            <label for="sizeRange" class="txt-primary">Data points</label>
          </div>
          <div class='col-md-8 pe-2'>
            <input type="range" class="form-range" min="2" max="501" step="1" id="sizeRange" v-model='frequency_points' @change='setCalibrated(false)'>
          </div>
          <div class='col-md-2'>
            <label class='txt-primary'>{{frequency_points}}</label>
          </div>
        </div>
      
    </div>

    
    <div class='d-flex flex-column'>
        <!-- CALIBRATION COMMANDS-->
        <div class="input-group mb-2" @mousedown="setDraggable(false)" @mouseup="setDraggable(true)">
          <span class="input-group-text txt-primary col-sm-3">Start</span>
          <input type="number" :class="(parseFloat(frequency) >= 0 && parseFloat(frequency) < maxFrequency) ? 'form-control' : 'form-control is-invalid'" :min='minFrequency' :max='maxFrequency' placeholder="Start frequency" aria-label="freq" aria-describedby="basic-addon1" id="freq" v-model="frequency" @change="setCalibrated(false)">
          <span class="input-group-text txt-primary" id="basic-addon1">MHz</span>
        </div>

        <div class="input-group" @mousedown="setDraggable(false)" @mouseup="setDraggable(true)">
          <span class="input-group-text txt-primary col-sm-3" id="basic-addon1">End</span>
          <input type="number" :class="(parseFloat(frequency_end) >= 0 && parseFloat(frequency_end) <= maxFrequency && parseFloat(frequency_end) > frequency) ? 'form-control' : 'form-control is-invalid'" :min='minFrequency' :max='maxFrequency' placeholder="End frequency" aria-label="freq" aria-describedby="basic-addon1" id="freq_end" v-model="frequency_end" @change="setCalibrated(false)">     
          <span class="input-group-text txt-primary" id="basic-addon1">MHz</span>
        </div>
    </div>


    <!-- Drag and Drop elements -->
    <drag-and-drop-components id='dragdropcomponents' header="Calibration Standards" :display='standards' :syncPorts='getSyncPorts' @port1change='updatePort1' @port2change='updatePort2'/>


    <!-- Scan and Save logic DIFFERENT FROM THE NON-SYNCED PORTS VERSION-->
    <div class='d-flex flex-row justify-content-center form-check-inline'>
          <label class='txt-primary txt-lg me-2'>Short</label>
          <input class="form-check-input mt-2 me-2" type="checkbox" value="" id="shortSyncCheck" :checked='getShortSaved' disabled>
          <label class='txt-primary txt-lg me-2'>Open</label>
          <input class="form-check-input mt-2 me-2" type="checkbox" value="" id="openSyncCheck" :checked='getOpenSaved' disabled>
          <label class='txt-primary txt-lg me-2'>Load</label>
          <input class="form-check-input mt-2 me-2" type="checkbox" value="" id="loadSyncCheck" :checked='getLoadSaved' disabled>
          <label class='txt-primary txt-lg me-2'>Thru</label>
          <input class="form-check-input mt-2 me-2" type="checkbox" value="" id="thruCheck" :checked='getThruSaved' disabled>
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
                  <p>This could take around 1 minute.</p>
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
   name: 'DragAndDropCalibrateSyncPorts',
    props: ['sparams'],
    components:{
        DragAndDropComponents,
    },
    data () {
        return {
            frequency_points: 20,
            frequency: 1.0,   //MHz
            frequency_end: 4.0, //MHz
            minFrequency: 0,
            maxFrequency: 4000.0, //MHz
            units: 1E6, //MHz
            avgCounts: 1,
            islog: false,
            s11: this.sparams.includes('s11') ? true:false,         
            s12: this.sparams.includes('s12') ? true:false,
            s21: this.sparams.includes('s21') ? true:false,
            s22: this.sparams.includes('s22') ? true:false,

            standards: [
                {type: 'short', img: "Short"},
                {type: 'open', img: "Open"},
                {type: 'load', img: "Load"},
                {type: 'thru', img: "Thru"},
            ],
            calibration: [
                {type: 'short', port:"sync", required: true, scanned: false, saved: false},
                {type: 'open', port:"sync", required: true, scanned: false, saved: false},
                {type: 'load', port:"sync", required: true, scanned: false, saved: false},
                {type: 'thru', port:"sync", required: true, scanned: false, saved: false},
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
        getShortSaved(){
            if(this.calibration[0].saved){
                return true;
            } else{
                return false;
            }
        },
        getOpenSaved(){
            if(this.calibration[1].saved){
                return true;
            } else{
                return false;
            }
        },
        getLoadSaved(){
            if(this.calibration[2].saved){
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
            if((this.port1 == 'short' || this.port2 == 'short') && this.calibration[0].scanned){
                return true;
            } else if((this.port1 == 'open' || this.port2 == 'open') && this.calibration[1].scanned){
                return true;
            } else if((this.port1 == 'load' || this.port2 == 'load') && this.calibration[2].scanned){
                return true;
            }  else if((this.port1 == 'thru' || this.port2 == 'thru') && this.calibration[3].scanned){
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
            this.calibration.forEach((connection) => {
                if(connection.type == this.port1 || connection.type == this.port2){
                    connection.scanned = true;
                }
            });
            
            // ACTUALLY SEND A COMMAND HERE TO RECEIVE CURRENT DATA
            let params = {
                t:dayjs().unix(),
                avg:Number(this.avgCounts),
                sparam:{s11:this.s11,s12:this.s12,s21:this.s21,s22:this.s22},
                what: this.port1 != '' ? this.port1 : this.port2  
            }

            this.$store.dispatch('requestRangeBeforeCal', params);
            this.$store.dispatch('setShowScanningModal', true);
        }, 
        save(){
            this.calibration.forEach((connection) => {
                if(connection.type == this.port1 || connection.type == this.port2){
                    if(connection.scanned){
                        connection.saved = true;
                    }
                }
            });
            
        },
        rangeFreqCalibration(){
            //command structure: {"cmd":"rc","range":{"Start":100000,"End":4000000},"size":2,"isLog":true,"avg":1,"sparam":{"S11":true,"S12":false,"S21":true,"S22":false}}
            if(this.frequency < this.minFrequency){
                this.frequency = this.minFrequency;
            } else if(this.frequency > this.maxFrequency){
                this.frequency = this.maxFrequency - 1;
            }

            if(this.frequency_end < this.frequency){
                this.frequency_end = this.frequency + 1;
            } else if(this.frequency_end > this.maxFrequency){
                this.frequency_end = this.maxFrequency;
            }

            
            let params = {
                id:'',
                t:dayjs().unix(),
                range:{start: this.frequency*this.units, end: this.frequency_end*this.units},
                size: Number(this.frequency_points),
                islog: this.islog,
                avg:Number(this.avgCounts),
                sparam:{s11:this.s11,s12:this.s12,s21:this.s21,s22:this.s22}
            }

            this.$store.dispatch('setCalibrationState', {sparams: params.sparam, average: Number(this.avgCounts), points: this.frequency_points, start: this.frequency, end: this.frequency_end});
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