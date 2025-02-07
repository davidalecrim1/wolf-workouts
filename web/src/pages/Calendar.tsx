import { NavBar } from "../components/NavBar"
import { PageHeader } from "../components/PageHeader"
import { Calendar as BigCalendar, momentLocalizer } from 'react-big-calendar'
import moment from 'moment'
import 'react-big-calendar/lib/css/react-big-calendar.css'

const localizer = momentLocalizer(moment)

type CalendarEvent = {
    id: number;
    title: string;
    start: Date;
    end: Date;
}

const dummyEventList: CalendarEvent[] = [
    {
        id: 1,
        title: 'Event 1',
        start: new Date('2025-02-07T09:00:00'),
        end: new Date('2025-02-07T10:00:00')
    },
    {
        id: 2,
        title: 'Event 2',
        start: new Date('2025-02-07T11:00:00'),
        end: new Date('2025-02-07T12:00:00')
    }
]

export const Calendar = () => {
    // There seems to be a type error with the BigCalendar component
    // This is a workaround to fix the type error
    const MyCalendar = BigCalendar as any;

    return (
        <div>
            <NavBar />
            <PageHeader title="Trainer's schedule" />
            <div className="container mx-auto px-4">
                <MyCalendar
                    localizer={localizer}
                    events={dummyEventList}
                    startAccessor="start"
                    endAccessor="end"
                    style={{
                        height: 800
                    }}
                />
            </div>
        </div>
    )
}