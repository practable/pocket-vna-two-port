<template>
  
 <div class='container-sm m-2 background-white border rounded'>
    
    <div v-if='!getParametersSet'>
         <h4 class='txt-primary txt-lg'>Please set the calibration parameters before progressing</h4>
     </div>

     <div v-else>
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
            <button id="request_calibration" type='button' class="button-lg button-tertiary" @click="confirmCal" :disabled='!ready_to_calibrate'>Calibrate</button>
        </div>

        



        <transition name='fade'>
        <div v-if='getShowCalibrationModal' class="modal" id='modal-show' tabindex="-1">
            <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header" style="background-color: #ccc">
                <h5 class="modal-title">Pocket VNA calibration complete</h5>
                </div>
                <div class="modal-body">
                <div class='d-flex row align-items-center'>
                    <div class='col-2' id="reveal-tick">
                        <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" fill="green" class="bi bi-check-lg" viewBox="0 0 16 16">
                            <path d="M12.736 3.97a.733.733 0 0 1 1.047 0c.286.289.29.756.01 1.05L7.88 12.01a.733.733 0 0 1-1.065.02L3.217 8.384a.757.757 0 0 1 0-1.06.733.733 0 0 1 1.047 0l3.052 3.093 5.4-6.425z"/>
                        </svg>
                    </div>
                    <div class='col-10'>
                    <p> You have correctly calibrated the Pocket VNA.</p>
                    <p>Please move on to the next tab, 3) Verification</p>
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

        <transition name='fade'>
        <div v-if='getShowErrorModal' class="modal" id='modal-show' tabindex="-1">
            <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header" style="background-color: #ccc">
                <h5 class="modal-title">Calibration incomplete</h5>
                </div>
                <div class="modal-body">
                <div class='d-flex row align-items-center'>
                    <div class='col-2' id="reveal-tick">
                        <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" fill="red" class="bi bi-exclamation-octagon" viewBox="0 0 16 16">
                            <path d="M4.54.146A.5.5 0 0 1 4.893 0h6.214a.5.5 0 0 1 .353.146l4.394 4.394a.5.5 0 0 1 .146.353v6.214a.5.5 0 0 1-.146.353l-4.394 4.394a.5.5 0 0 1-.353.146H4.893a.5.5 0 0 1-.353-.146L.146 11.46A.5.5 0 0 1 0 11.107V4.893a.5.5 0 0 1 .146-.353zM5.1 1 1 5.1v5.8L5.1 15h5.8l4.1-4.1V5.1L10.9 1z"/>
                            <path d="M7.002 11a1 1 0 1 1 2 0 1 1 0 0 1-2 0M7.1 4.995a.905.905 0 1 1 1.8 0l-.35 3.507a.552.552 0 0 1-1.1 0z"/>
                        </svg>
                    </div>
                    <div class='col-10'>
                    <p> The calibration process is not complete, see the message below:</p>
                    <p>{{ getErrorMessage }}</p>
                    </div>
                </div>
                </div>
                <div class="modal-footer">
                <button type="button" class="btn btn-danger" @click="setShowErrorModal(false)">Close</button>
                </div>
            </div>
            </div>
        </div>
        </transition>

    </div>

</div>


</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import DragAndDropComponents from '@/components/DragAndDropComponents.vue';
import dayjs from "dayjs";

export default {
   name: 'DragAndDropCalibrateSyncPorts',
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
            port1: '',
            port2: '',

        }
    },
    computed:{
        ...mapGetters([
            'getCalibrated',
            'getShowCalibrationModal',
            'getSyncPorts',
            'getShowScanningModal',
            'getParametersSet',
            'getCalibrationPorts',
            'getShowErrorModal',
            'getErrorMessage'
        ]),
        ready_to_calibrate(){
            let ready = true;
            this.getCalibrationPorts.forEach((connection) => {
                if(connection.required){
                    if(!connection.saved){
                        ready = false;
                    } 
                }
            });
            return ready;
        },
        getShortSaved(){
            if(this.getCalibrationPorts[0].saved){
                return true;
            } else{
                return false;
            }
        },
        getOpenSaved(){
            if(this.getCalibrationPorts[1].saved){
                return true;
            } else{
                return false;
            }
        },
        getLoadSaved(){
            if(this.getCalibrationPorts[2].saved){
                return true;
            } else{
                return false;
            }
        },
        getThruSaved(){
            if(this.getCalibrationPorts[3].saved){
                return true;
            } else{
                return false;
            }
        },
        getShowSave(){
            if((this.port1 == 'short' || this.port2 == 'short') && this.getCalibrationPorts[0].scanned){
                return true;
            } else if((this.port1 == 'open' || this.port2 == 'open') && this.getCalibrationPorts[1].scanned){
                return true;
            } else if((this.port1 == 'load' || this.port2 == 'load') && this.getCalibrationPorts[2].scanned){
                return true;
            }  else if((this.port1 == 'thru' || this.port2 == 'thru') && this.getCalibrationPorts[3].scanned){
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
            'setShowScanningModal',
            'setShowErrorModal'
        ]),
        scan(){
            let params = {
                t:dayjs().unix(),
                // avg:Number(this.calibrationState.avg),
                // sparam:this.calibrationState.sparam,
                what: this.port1 != '' ? this.port1 : this.port2  
            }

            this.$store.dispatch('requestRangeBeforeCal', params);
            this.$store.dispatch('setShowScanningModal', true);
        }, 
        save(){
            this.getCalibrationPorts.forEach((connection) => {
                if(connection.type == this.port1 || connection.type == this.port2){
                    if(connection.scanned){
                        connection.saved = true;
                    }
                }
            });
            
        },
        confirmCal(){
            let params = {
                t:dayjs().unix(),
                // avg:Number(this.calibrationState.avg),
                // sparam:this.calibrationState.sparam,
            }

            this.$store.dispatch('confirmCal', params);

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