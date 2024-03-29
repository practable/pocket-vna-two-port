<template>
  
 <div class='container-sm m-2 background-white border rounded'>

     <div v-if='!getCalibrated'>
         <h4 class='txt-primary txt-lg'>Please calibrate the Pocket VNA before progressing</h4>
     </div>

    <div v-else>
        
        <!-- Drag and Drop elements -->
        <drag-and-drop-components id='dragdropcomponents' header="Calibration Standards" :display='standards' :syncPorts="getSyncPorts" @port1change='updatePort1' @port2change='updatePort2'/>
        
        <div class='d-flex flex-row justify-content-center'>

            <button id="verify" type='button' class="button-lg button-primary" @click='verify' :disabled="port1 == '' && port2 == ''">Verify</button>

        </div>

        <div class='d-flex flex-row justify-content-center'>

            <button id="useCalibration" type='button' class="button-lg button-success" @click='useCalibration'>Use Calibration</button>

        </div>

        
        <transition name='fade'>
        <div v-if='getShowVerifiedModal' class="modal" id='modal-show' tabindex="-1">
            <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header" style="background-color: #ccc">
                <h5 class="modal-title">Pocket VNA calibration verified</h5>
                </div>
                <div class="modal-body">
                <div class='d-flex row align-items-center'>
                    <div class='col-2' id="reveal-tick">
                        <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" fill="green" class="bi bi-check-lg" viewBox="0 0 16 16">
                            <path d="M12.736 3.97a.733.733 0 0 1 1.047 0c.286.289.29.756.01 1.05L7.88 12.01a.733.733 0 0 1-1.065.02L3.217 8.384a.757.757 0 0 1 0-1.06.733.733 0 0 1 1.047 0l3.052 3.093 5.4-6.425z"/>
                        </svg>
                    </div>
                    <div class='col-10'>
                    <p> You have verified the Pocket VNA calibration.</p>
                    <p>Please move on to the next tab, 4) Measurement</p>
                    </div>
                </div>
                </div>
                <div class="modal-footer">
                <button type="button" class="btn btn-danger" @click="setShowVerifiedModal(false)">Close</button>
                </div>
            </div>
            </div>
        </div>
        </transition>


        <transition name='fade'>
        <div v-if='getShowRequestModal' class="modal" id='modal-show' tabindex="-1">
            <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header" style="background-color: #ccc">
                <h5 class="modal-title">Requesting Measurement from Pocket VNA</h5>
                </div>
                <div class="modal-body">
                <div class='d-flex row align-items-center'>
                    <div class='col-2'>
                    <div class="spinner-border text-primary text-center" role="status">
                        <span class="visually-hidden">Loading...</span>
                    </div>
                    </div>
                    <div class='col-10'>
                    <p> Please wait for the pocket VNA to measure.</p>
                    <p>This could take around 30 seconds.</p>
                    </div>
                </div>
                </div>
                <div class="modal-footer">
                <button type="button" class="btn btn-danger" @click="setShowRequestModal(false)">Close</button>
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
   name: 'DragAndDropVerification',
    props: ['calibrationState'],
    components:{
        DragAndDropComponents,
    },
    data () {
        return {
            standards: [
                {type: 'short', img: "Short"},
                {type: 'open', img: "Open"},
                {type: 'load', img: "Load"},
                {type: 'thru', img: "Thru"}
            ],

            port1: '',
            port2: '',

        }
    },
    computed:{
        ...mapGetters([
            'getCalibrated',
            'getShowRequestModal',
            'getShowVerifiedModal',
            'getSyncPorts'
        ]),
        
    },
    watch:{
        
    },
    created(){
        
    },
    mounted(){
        

    },
    methods:{
        ...mapActions([
            'setShowRequestModal',
            'setShowVerifiedModal'
        ]),
        // Take a measurement of calibration standards to check the calibration process
        verify(){
            this.rangeFreqRequest();
        }, 
        //Set the state as calibrated to allow measurements to be taken.
        useCalibration(){
            this.$store.dispatch('setVerified', true);
            this.$store.dispatch('setShowVerifiedModal', true);
        },
      rangeFreqRequest(){
        //command structure: {"cmd":"crq","avg":1,"sparam":{"S11":true,"S12":false,"S21":false,"S22":false}}
        let params = {
            t:dayjs().unix(),
            avg:Number(this.calibrationState.avg),
            sparam:this.calibrationState.sparam,
            what: this.port1 != '' ? this.port1 : this.port2  
          }

          this.$store.dispatch('requestRangeAfterCal', params);
          this.$store.dispatch('setShowRequestModal', true);
          
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