import styles from '../styles/Index.module.css'
import { useRef } from 'react'

export default function Home() {

  const submitForm = async ()=>{
    const formData = new FormData()
    let inputFile = document.querySelector('input[type="file"]').files[0]
    let questions = document.getElementById("questions").value
    let answers = document.getElementById("answers").value
    let comments = document.getElementById("comments").value

    console.log(inputFile)

    // Append the data to the form
    formData.append("file", inputFile)
    formData.append("questions", questions)
    formData.append("answers", answers)
    formData.append("comments", comments)

    const resp = await fetch('http://localhost:5000/parsefile',{
      method: "POST",
      body: formData
    })

    console.log(resp)
  }

  return (
    <div className={styles.container}>
      <input type='file' id='file' /> <br />
      <label htmlFor='questions'>Column Num with Questions: </label>
      <input type='number' id='questions' /><br />
      <label htmlFor='answers'>Column Num with Answers: </label>
      <input type='number' id='answers' /><br />
      <label htmlFor='comments'>Column Num with Comments: </label>
      <input type='number' id='comments' /><br />
      <button onClick={submitForm}>Upload File</button>

    </div>
  )
}
