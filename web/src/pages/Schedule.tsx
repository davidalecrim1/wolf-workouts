import { NavBar } from "../components/NavBar"
import { PageHeader } from "../components/PageHeader"

export const Schedule = () => {
    return (
        <div>
            <NavBar />
            <PageHeader title="Set schedule" />
            <div>Maybe I can reuse the calendar component here? It depends on the business logic</div>
        </div>
    )
}
