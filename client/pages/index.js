import styles from '../styles/Index.module.css'
import { useRef } from 'react'

export default function Home() {

  const inputFile = useRef(null)

  const submitForm = async ()=>{
    const formData = new FormData()
    var inputFile = document.querySelector('input[type="file"]').files[0]

    console.log(inputFile)

    formData.append("file", inputFile)

    const resp = await fetch('http://localhost:5000/parsefile',{
      method: "POST",
      body: formData
    })

    console.log(resp)
  }

  return (
    <div className={styles.container}>
      <form encType='multipart/form-data'>
        <input type='file' id='file' ref={inputFile} />
      </form>
      <button onClick={submitForm}>Upload File</button>

    </div>
  )
}
