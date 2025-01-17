func or(channels ...<-chan interface{}) <-chan interface{} {
    // Создаем выходной канал
    out := make(chan interface{})
    
    // Запускаем горутину для обработки всех входных каналов
    go func() {
        defer close(out)
        
        // Используем WaitGroup для отслеживания всех горутин
        var wg sync.WaitGroup
        
        // Создаем функцию для прослушивания одного канала
        listen := func(ch <-chan interface{}) {
            defer wg.Done()
            for v := range ch {
                select {
                case out <- v:
                case <-out:
                    return
                }
            }
        }
        
        // Запускаем горутину для каждого входного канала
        wg.Add(len(channels))
        for _, ch := range channels {
            go listen(ch)
        }
        
        // Ждем завершения всех горутин
        wg.Wait()
    }()
    
    return out
}